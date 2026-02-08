#!/usr/bin/env node

"use strict";

const os = require("os");
const fs = require("fs");
const path = require("path");
const https = require("https");
const zlib = require("zlib");
const { spawn } = require("child_process");

const pkg = require("../package.json");
const VERSION = pkg.version;
const REPO = "atlantic-blue/greenlight";

const ARCH_MAP = { x64: "amd64", arm64: "arm64" };
const PLATFORM_MAP = { darwin: "darwin", linux: "linux" };

function getDownloadUrl() {
  const platform = PLATFORM_MAP[os.platform()];
  const arch = ARCH_MAP[os.arch()];

  if (!platform) {
    console.error(
      `Unsupported platform: ${os.platform()}. Greenlight supports macOS and Linux.`
    );
    process.exit(1);
  }

  if (!arch) {
    console.error(
      `Unsupported architecture: ${os.arch()}. Greenlight supports x64 and arm64.`
    );
    process.exit(1);
  }

  const filename = `greenlight_${VERSION}_${platform}_${arch}.tar.gz`;
  return `https://github.com/${REPO}/releases/download/v${VERSION}/${filename}`;
}

function fetch(url) {
  return new Promise((resolve, reject) => {
    https
      .get(url, (res) => {
        if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
          return fetch(res.headers.location).then(resolve, reject);
        }
        if (res.statusCode !== 200) {
          reject(
            new Error(
              `Download failed (HTTP ${res.statusCode}). ` +
                `Make sure v${VERSION} exists at https://github.com/${REPO}/releases`
            )
          );
          return;
        }
        resolve(res);
      })
      .on("error", reject);
  });
}

function extractTarGz(stream, destDir) {
  return new Promise((resolve, reject) => {
    const gunzip = zlib.createGunzip();
    const chunks = [];
    let totalSize = 0;

    stream
      .pipe(gunzip)
      .on("data", (chunk) => {
        chunks.push(chunk);
        totalSize += chunk.length;
      })
      .on("end", () => {
        const buffer = Buffer.concat(chunks, totalSize);
        let offset = 0;

        while (offset < buffer.length) {
          // tar header is 512 bytes
          if (offset + 512 > buffer.length) break;
          const header = buffer.subarray(offset, offset + 512);

          // Check for end-of-archive (two consecutive zero blocks)
          if (header.every((b) => b === 0)) break;

          const name = header.subarray(0, 100).toString("utf8").replace(/\0/g, "");
          const sizeOctal = header.subarray(124, 136).toString("utf8").replace(/\0/g, "").trim();
          const size = parseInt(sizeOctal, 8) || 0;
          const typeFlag = header[156];

          offset += 512; // move past header

          if (typeFlag === 48 || typeFlag === 0) {
            // Regular file (type '0' or null)
            const content = buffer.subarray(offset, offset + size);
            const filePath = path.join(destDir, path.basename(name));
            fs.writeFileSync(filePath, content);
            fs.chmodSync(filePath, 0o755);
          }

          // Advance past file data, padded to 512-byte boundary
          offset += Math.ceil(size / 512) * 512;
        }

        resolve();
      })
      .on("error", reject);
  });
}

async function run() {
  const url = getDownloadUrl();
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), "greenlight-"));
  const binaryPath = path.join(tempDir, "greenlight");

  try {
    const stream = await fetch(url);
    await extractTarGz(stream, tempDir);

    if (!fs.existsSync(binaryPath)) {
      console.error("Failed to extract greenlight binary from archive.");
      process.exit(1);
    }

    const args = process.argv.slice(2);
    const child = spawn(binaryPath, args, { stdio: "inherit" });

    child.on("close", (code) => {
      // Clean up temp directory
      fs.rmSync(tempDir, { recursive: true, force: true });
      process.exit(code);
    });

    child.on("error", (err) => {
      console.error(`Failed to run greenlight: ${err.message}`);
      fs.rmSync(tempDir, { recursive: true, force: true });
      process.exit(1);
    });
  } catch (err) {
    console.error(err.message);
    fs.rmSync(tempDir, { recursive: true, force: true });
    process.exit(1);
  }
}

run();
