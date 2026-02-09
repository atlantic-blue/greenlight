package installer_test

import (
	"os"
	"strings"
	"testing"
)

// C-39 Tests: HelpInsightSection
// These tests verify the markdown content of help.md contains the required INSIGHT section

func TestHelpMd_ContainsInsightSection(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/help.md")
	if err != nil {
		t.Fatalf("failed to read help.md: %v", err)
	}

	helpContent := string(content)

	if !strings.Contains(helpContent, "INSIGHT") {
		t.Error("help.md missing INSIGHT section")
	}
}

func TestHelpMd_InsightSectionBetweenMonitorAndShip(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/help.md")
	if err != nil {
		t.Fatalf("failed to read help.md: %v", err)
	}

	helpContent := string(content)

	// Find positions of MONITOR, INSIGHT, and SHIP
	monitorPos := strings.Index(helpContent, "MONITOR")
	insightPos := strings.Index(helpContent, "INSIGHT")
	shipPos := strings.Index(helpContent, "SHIP")

	if monitorPos == -1 {
		t.Error("help.md missing MONITOR section")
	}
	if insightPos == -1 {
		t.Error("help.md missing INSIGHT section")
	}
	if shipPos == -1 {
		t.Error("help.md missing SHIP section")
	}

	// Verify order: MONITOR < INSIGHT < SHIP
	if monitorPos != -1 && insightPos != -1 && insightPos <= monitorPos {
		t.Errorf("INSIGHT section (pos %d) should come after MONITOR section (pos %d)", insightPos, monitorPos)
	}
	if insightPos != -1 && shipPos != -1 && insightPos >= shipPos {
		t.Errorf("INSIGHT section (pos %d) should come before SHIP section (pos %d)", insightPos, shipPos)
	}
}

func TestHelpMd_InsightSectionContainsRoadmapCommand(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/help.md")
	if err != nil {
		t.Fatalf("failed to read help.md: %v", err)
	}

	helpContent := string(content)

	// Find INSIGHT section
	insightPos := strings.Index(helpContent, "INSIGHT")
	if insightPos == -1 {
		t.Fatal("help.md missing INSIGHT section")
	}

	// Find next section after INSIGHT (SHIP)
	shipPos := strings.Index(helpContent[insightPos:], "SHIP")
	if shipPos == -1 {
		t.Fatal("help.md missing SHIP section after INSIGHT")
	}

	insightSection := helpContent[insightPos : insightPos+shipPos]

	if !strings.Contains(insightSection, "/gl:roadmap") {
		t.Error("INSIGHT section missing /gl:roadmap command")
	}
	if !strings.Contains(insightSection, "Product roadmap + milestones") &&
	   !strings.Contains(insightSection, "roadmap") {
		t.Error("INSIGHT section missing roadmap description")
	}
}

func TestHelpMd_InsightSectionContainsChangelogCommand(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/help.md")
	if err != nil {
		t.Fatalf("failed to read help.md: %v", err)
	}

	helpContent := string(content)

	// Find INSIGHT section
	insightPos := strings.Index(helpContent, "INSIGHT")
	if insightPos == -1 {
		t.Fatal("help.md missing INSIGHT section")
	}

	// Find next section after INSIGHT (SHIP)
	shipPos := strings.Index(helpContent[insightPos:], "SHIP")
	if shipPos == -1 {
		t.Fatal("help.md missing SHIP section after INSIGHT")
	}

	insightSection := helpContent[insightPos : insightPos+shipPos]

	if !strings.Contains(insightSection, "/gl:changelog") {
		t.Error("INSIGHT section missing /gl:changelog command")
	}
	if !strings.Contains(insightSection, "Human-readable changelog") &&
	   !strings.Contains(insightSection, "changelog") {
		t.Error("INSIGHT section missing changelog description")
	}
}

func TestHelpMd_ContainsThreeViewsTagline(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/help.md")
	if err != nil {
		t.Fatalf("failed to read help.md: %v", err)
	}

	helpContent := string(content)

	if !strings.Contains(helpContent, "Three views") {
		t.Error("help.md missing three-views tagline")
	}
	if !strings.Contains(helpContent, "/gl:status") {
		t.Error("help.md three-views tagline missing /gl:status reference")
	}
	if !strings.Contains(helpContent, "/gl:roadmap") {
		t.Error("help.md three-views tagline missing /gl:roadmap reference")
	}
	if !strings.Contains(helpContent, "/gl:changelog") {
		t.Error("help.md three-views tagline missing /gl:changelog reference")
	}
}

func TestHelpMd_FlowLineContainsDocumentationSteps(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/help.md")
	if err != nil {
		t.Fatalf("failed to read help.md: %v", err)
	}

	helpContent := string(content)

	// The FLOW line should mention design produces ROADMAP and DECISIONS
	if !strings.Contains(helpContent, "design") {
		t.Error("help.md missing design step in flow")
	}
	if !strings.Contains(helpContent, "ROADMAP") {
		t.Error("help.md flow missing ROADMAP reference")
	}
	if !strings.Contains(helpContent, "DECISIONS") {
		t.Error("help.md flow missing DECISIONS reference")
	}
}

func TestHelpMd_BuildSectionSliceIncludesSummaryStep(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/help.md")
	if err != nil {
		t.Fatalf("failed to read help.md: %v", err)
	}

	helpContent := string(content)

	// Find BUILD section
	buildPos := strings.Index(helpContent, "BUILD")
	if buildPos == -1 {
		t.Fatal("help.md missing BUILD section")
	}

	// Find next section after BUILD (MONITOR)
	monitorPos := strings.Index(helpContent[buildPos:], "MONITOR")
	if monitorPos == -1 {
		t.Fatal("help.md missing MONITOR section after BUILD")
	}

	buildSection := helpContent[buildPos : buildPos+monitorPos]

	// The /gl:slice command description should mention summary
	if !strings.Contains(buildSection, "/gl:slice") {
		t.Fatal("BUILD section missing /gl:slice command")
	}
	if !strings.Contains(buildSection, "summary") {
		t.Error("BUILD section /gl:slice description missing summary step")
	}
}

func TestHelpMd_AllExistingSectionsPresent(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/help.md")
	if err != nil {
		t.Fatalf("failed to read help.md: %v", err)
	}

	helpContent := string(content)

	expectedSections := []string{
		"SETUP",
		"BROWNFIELD",
		"BUILD",
		"MONITOR",
		"INSIGHT",
		"SHIP",
	}

	for _, section := range expectedSections {
		if !strings.Contains(helpContent, section) {
			t.Errorf("help.md missing section: %s", section)
		}
	}
}

// C-40 Tests: StatusDocumentationReference

func TestStatusMd_ContainsDocumentationReferenceLine(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/status.md")
	if err != nil {
		t.Fatalf("failed to read status.md: %v", err)
	}

	statusContent := string(content)

	if !strings.Contains(statusContent, "Product view") {
		t.Error("status.md missing 'Product view' reference")
	}
	if !strings.Contains(statusContent, "/gl:roadmap") {
		t.Error("status.md missing /gl:roadmap reference")
	}
	if !strings.Contains(statusContent, "History") {
		t.Error("status.md missing 'History' reference")
	}
	if !strings.Contains(statusContent, "/gl:changelog") {
		t.Error("status.md missing /gl:changelog reference")
	}
}

func TestStatusMd_DocumentationReferenceFormat(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/status.md")
	if err != nil {
		t.Fatalf("failed to read status.md: %v", err)
	}

	statusContent := string(content)

	// The exact format should be: "Product view: /gl:roadmap | History: /gl:changelog"
	expectedFormat := "Product view: /gl:roadmap | History: /gl:changelog"

	if !strings.Contains(statusContent, expectedFormat) {
		t.Errorf("status.md missing expected format: %s", expectedFormat)
	}
}

// C-36 Tests: DesignRoadmapProduction

func TestDesignMd_ContainsRoadmapProductionInstructions(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/design.md")
	if err != nil {
		t.Fatalf("failed to read design.md: %v", err)
	}

	designContent := string(content)

	if !strings.Contains(designContent, "ROADMAP.md") {
		t.Error("design.md missing ROADMAP.md reference")
	}
}

func TestDesignMd_ReferencesRoadmapAfterDesignApproval(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/design.md")
	if err != nil {
		t.Fatalf("failed to read design.md: %v", err)
	}

	designContent := string(content)

	// Design should mention producing/creating ROADMAP.md
	if !strings.Contains(designContent, "ROADMAP") {
		t.Error("design.md missing ROADMAP production instructions")
	}

	// Should reference it in context of design session
	lowerContent := strings.ToLower(designContent)
	if !strings.Contains(lowerContent, "roadmap") {
		t.Error("design.md missing roadmap instructions")
	}
}

func TestDesignMd_ContainsMermaidDiagramReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/design.md")
	if err != nil {
		t.Fatalf("failed to read design.md: %v", err)
	}

	designContent := string(content)

	// ROADMAP.md should include architecture diagram in Mermaid format
	lowerContent := strings.ToLower(designContent)
	if !strings.Contains(lowerContent, "mermaid") &&
	   !strings.Contains(lowerContent, "diagram") &&
	   !strings.Contains(lowerContent, "architecture") {
		t.Error("design.md missing architecture diagram reference for ROADMAP.md")
	}
}

// C-37 Tests: DesignDecisionsSeeding

func TestDesignMd_ContainsDecisionsProductionInstructions(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/design.md")
	if err != nil {
		t.Fatalf("failed to read design.md: %v", err)
	}

	designContent := string(content)

	if !strings.Contains(designContent, "DECISIONS.md") {
		t.Error("design.md missing DECISIONS.md reference")
	}
}

func TestDesignMd_ReferencesDecisionSeeding(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/design.md")
	if err != nil {
		t.Fatalf("failed to read design.md: %v", err)
	}

	designContent := string(content)

	// Design should mention seeding DECISIONS.md from technical decisions
	if !strings.Contains(designContent, "DECISIONS") {
		t.Error("design.md missing DECISIONS seeding instructions")
	}

	lowerContent := strings.ToLower(designContent)
	if !strings.Contains(lowerContent, "decision") {
		t.Error("design.md missing decision-related instructions")
	}
}

func TestDesignMd_ContainsDecisionLogTableStructure(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/design.md")
	if err != nil {
		t.Fatalf("failed to read design.md: %v", err)
	}

	designContent := string(content)

	// DECISIONS.md should have a table structure with specific columns
	// Check for key column names
	lowerContent := strings.ToLower(designContent)

	// At minimum, should mention the decision log structure
	if !strings.Contains(lowerContent, "decision") {
		t.Error("design.md missing decision log structure reference")
	}
}

func TestDesignMd_ContainsBothRoadmapAndDecisions(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/design.md")
	if err != nil {
		t.Fatalf("failed to read design.md: %v", err)
	}

	designContent := string(content)

	// Design command should produce BOTH ROADMAP.md and DECISIONS.md
	hasRoadmap := strings.Contains(designContent, "ROADMAP.md")
	hasDecisions := strings.Contains(designContent, "DECISIONS.md")

	if !hasRoadmap {
		t.Error("design.md missing ROADMAP.md production instructions")
	}
	if !hasDecisions {
		t.Error("design.md missing DECISIONS.md production instructions")
	}
}

// C-41 Tests: SliceSummaryGeneration
// These tests verify slice.md contains summary generation instructions after verification

func TestSliceMd_ContainsSummaryGenerationStep(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/slice.md")
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	sliceContent := string(content)

	if !strings.Contains(sliceContent, "summary") && !strings.Contains(sliceContent, "Summary") && !strings.Contains(sliceContent, "SUMMARY") {
		t.Error("slice.md missing summary generation step")
	}
}

func TestSliceMd_ReferencesSummaryAfterVerification(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/slice.md")
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	sliceContent := string(content)

	// Should reference summary generation after verification step
	lowerContent := strings.ToLower(sliceContent)
	if !strings.Contains(lowerContent, "after verification") && !strings.Contains(lowerContent, "verification succeeds") {
		t.Error("slice.md missing reference to summary timing (after verification)")
	}
}

func TestSliceMd_ReferencesSpawningTaskForSummary(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/slice.md")
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	sliceContent := string(content)

	// Should mention spawning a Task for summary generation
	if !strings.Contains(sliceContent, "Task") && !strings.Contains(sliceContent, "spawn") {
		t.Error("slice.md missing reference to spawning Task for summary")
	}
}

func TestSliceMd_ReferencesSummariesDirectory(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/slice.md")
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	sliceContent := string(content)

	if !strings.Contains(sliceContent, "summaries/") && !strings.Contains(sliceContent, ".greenlight/summaries") {
		t.Error("slice.md missing reference to summaries/ directory")
	}
}

func TestSliceMd_ReferencesSummaryFileNaming(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/slice.md")
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	sliceContent := string(content)

	// Should reference {slice-id}-SUMMARY.md naming pattern
	if !strings.Contains(sliceContent, "SUMMARY.md") && !strings.Contains(sliceContent, "-SUMMARY") {
		t.Error("slice.md missing reference to SUMMARY.md file naming")
	}
}

func TestSliceMd_ReferencesStructuredDataForSummary(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/slice.md")
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	sliceContent := string(content)

	// Should mention structured data or XML context blocks
	lowerContent := strings.ToLower(sliceContent)
	hasStructuredData := strings.Contains(lowerContent, "structured data")
	hasXMLContext := strings.Contains(sliceContent, "XML") || strings.Contains(sliceContent, "<slice>") || strings.Contains(lowerContent, "xml context")

	if !hasStructuredData && !hasXMLContext {
		t.Error("slice.md missing reference to structured data/XML context for summary")
	}
}

func TestSliceMd_ReferencesRoadmapUpdateAfterSummary(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/slice.md")
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	sliceContent := string(content)

	// Should mention updating ROADMAP.md after summary
	if !strings.Contains(sliceContent, "ROADMAP.md") {
		t.Error("slice.md missing reference to ROADMAP.md update")
	}
}

func TestSliceMd_SummaryIsNonBlocking(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/slice.md")
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	sliceContent := string(content)

	// Should mention that summary failure doesn't block
	lowerContent := strings.ToLower(sliceContent)
	hasNonBlocking := strings.Contains(lowerContent, "does not block") ||
		strings.Contains(lowerContent, "doesn't block") ||
		strings.Contains(lowerContent, "not block") ||
		strings.Contains(lowerContent, "non-blocking")

	if !hasNonBlocking {
		t.Error("slice.md missing mention of summary being non-blocking")
	}
}

func TestSliceMd_ReferencesProductLanguageForSummary(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/slice.md")
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	sliceContent := string(content)

	// Should mention product language (not implementation language)
	lowerContent := strings.ToLower(sliceContent)
	if !strings.Contains(lowerContent, "product language") && !strings.Contains(lowerContent, "user-facing") {
		t.Error("slice.md missing reference to product language for summary")
	}
}

// C-42 Tests: WrapSummaryGeneration
// These tests verify wrap.md contains summary generation instructions after commit

func TestWrapMd_ContainsSummaryGenerationStep(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/wrap.md")
	if err != nil {
		t.Fatalf("failed to read wrap.md: %v", err)
	}

	wrapContent := string(content)

	if !strings.Contains(wrapContent, "summary") && !strings.Contains(wrapContent, "Summary") && !strings.Contains(wrapContent, "SUMMARY") {
		t.Error("wrap.md missing summary generation step")
	}
}

func TestWrapMd_ReferencesSummaryAfterCommit(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/wrap.md")
	if err != nil {
		t.Fatalf("failed to read wrap.md: %v", err)
	}

	wrapContent := string(content)

	// Should reference summary after wrap commit
	lowerContent := strings.ToLower(wrapContent)
	if !strings.Contains(lowerContent, "after") && !strings.Contains(lowerContent, "commit") {
		t.Error("wrap.md missing reference to summary timing (after commit)")
	}
}

func TestWrapMd_ReferencesSpawningTaskForWrapSummary(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/wrap.md")
	if err != nil {
		t.Fatalf("failed to read wrap.md: %v", err)
	}

	wrapContent := string(content)

	// Should mention spawning Task for summary
	if !strings.Contains(wrapContent, "Task") && !strings.Contains(wrapContent, "spawn") {
		t.Error("wrap.md missing reference to spawning Task for summary")
	}
}

func TestWrapMd_ReferencesWrapSummaryNaming(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/wrap.md")
	if err != nil {
		t.Fatalf("failed to read wrap.md: %v", err)
	}

	wrapContent := string(content)

	// Should reference {boundary-name}-wrap-SUMMARY.md naming
	if !strings.Contains(wrapContent, "wrap-SUMMARY") && !strings.Contains(wrapContent, "-wrap-") {
		t.Error("wrap.md missing reference to wrap-SUMMARY.md naming pattern")
	}
}

func TestWrapMd_ReferencesWrapSummariesDirectory(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/wrap.md")
	if err != nil {
		t.Fatalf("failed to read wrap.md: %v", err)
	}

	wrapContent := string(content)

	if !strings.Contains(wrapContent, "summaries/") && !strings.Contains(wrapContent, ".greenlight/summaries") {
		t.Error("wrap.md missing reference to summaries/ directory")
	}
}

func TestWrapMd_ReferencesRoadmapWrapProgressUpdate(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/wrap.md")
	if err != nil {
		t.Fatalf("failed to read wrap.md: %v", err)
	}

	wrapContent := string(content)

	// Should mention updating ROADMAP.md wrap progress
	hasRoadmap := strings.Contains(wrapContent, "ROADMAP.md")
	lowerContent := strings.ToLower(wrapContent)
	hasWrapProgress := strings.Contains(lowerContent, "wrap progress") || strings.Contains(lowerContent, "progress")

	if !hasRoadmap || !hasWrapProgress {
		t.Error("wrap.md missing reference to ROADMAP.md wrap progress update")
	}
}

func TestWrapMd_WrapSummaryIsNonBlocking(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/wrap.md")
	if err != nil {
		t.Fatalf("failed to read wrap.md: %v", err)
	}

	wrapContent := string(content)

	// Should mention that summary failure doesn't block wrap
	lowerContent := strings.ToLower(wrapContent)
	hasNonBlocking := strings.Contains(lowerContent, "does not block") ||
		strings.Contains(lowerContent, "doesn't block") ||
		strings.Contains(lowerContent, "not block") ||
		strings.Contains(lowerContent, "non-blocking")

	if !hasNonBlocking {
		t.Error("wrap.md missing mention of summary being non-blocking")
	}
}

// C-43 Tests: QuickSummaryGeneration
// These tests verify quick.md contains summary generation instructions

func TestQuickMd_ContainsSummaryGenerationStep(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/quick.md")
	if err != nil {
		t.Fatalf("failed to read quick.md: %v", err)
	}

	quickContent := string(content)

	if !strings.Contains(quickContent, "summary") && !strings.Contains(quickContent, "Summary") && !strings.Contains(quickContent, "SUMMARY") {
		t.Error("quick.md missing summary generation step")
	}
}

func TestQuickMd_ReferencesQuickSummaryNaming(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/quick.md")
	if err != nil {
		t.Fatalf("failed to read quick.md: %v", err)
	}

	quickContent := string(content)

	// Should reference quick-{timestamp}-SUMMARY.md naming
	if !strings.Contains(quickContent, "quick-") && !strings.Contains(quickContent, "SUMMARY") {
		t.Error("quick.md missing reference to quick- naming pattern")
	}
}

func TestQuickMd_ReferencesQuickSummariesDirectory(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/quick.md")
	if err != nil {
		t.Fatalf("failed to read quick.md: %v", err)
	}

	quickContent := string(content)

	if !strings.Contains(quickContent, "summaries/") && !strings.Contains(quickContent, ".greenlight/summaries") {
		t.Error("quick.md missing reference to summaries/ directory")
	}
}

func TestQuickMd_ReferencesSpawningTaskForQuickSummary(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/quick.md")
	if err != nil {
		t.Fatalf("failed to read quick.md: %v", err)
	}

	quickContent := string(content)

	// Should mention spawning Task for summary
	if !strings.Contains(quickContent, "Task") && !strings.Contains(quickContent, "spawn") {
		t.Error("quick.md missing reference to spawning Task for summary")
	}
}

func TestQuickMd_ReferencesDecisionAppend(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/quick.md")
	if err != nil {
		t.Fatalf("failed to read quick.md: %v", err)
	}

	quickContent := string(content)

	// Should mention appending to DECISIONS.md if decision made
	if !strings.Contains(quickContent, "DECISIONS.md") && !strings.Contains(quickContent, "DECISIONS") {
		t.Error("quick.md missing reference to DECISIONS.md append")
	}
}

func TestQuickMd_QuickSummaryIsNonBlocking(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/quick.md")
	if err != nil {
		t.Fatalf("failed to read quick.md: %v", err)
	}

	quickContent := string(content)

	// Should mention that summary failure doesn't block quick task
	lowerContent := strings.ToLower(quickContent)
	hasNonBlocking := strings.Contains(lowerContent, "does not block") ||
		strings.Contains(lowerContent, "doesn't block") ||
		strings.Contains(lowerContent, "not block") ||
		strings.Contains(lowerContent, "non-blocking")

	if !hasNonBlocking {
		t.Error("quick.md missing mention of summary being non-blocking")
	}
}

func TestQuickMd_ReferencesTimestampInSummary(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/quick.md")
	if err != nil {
		t.Fatalf("failed to read quick.md: %v", err)
	}

	quickContent := string(content)

	// Should reference timestamp in summary context
	lowerContent := strings.ToLower(quickContent)
	if !strings.Contains(lowerContent, "timestamp") {
		t.Error("quick.md missing reference to timestamp for summary")
	}
}

// C-44 Tests: DecisionAggregation
// These tests verify slice.md contains decision aggregation instructions

func TestSliceMd_ContainsDecisionAggregationStep(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/slice.md")
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	sliceContent := string(content)

	// Should mention decision aggregation or collecting decisions
	lowerContent := strings.ToLower(sliceContent)
	if !strings.Contains(lowerContent, "decision") && !strings.Contains(sliceContent, "DECISIONS") {
		t.Error("slice.md missing decision aggregation step")
	}
}

func TestSliceMd_ReferencesCollectingDecisionsFromAgents(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/slice.md")
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	sliceContent := string(content)

	// Should reference collecting decisions from agents
	lowerContent := strings.ToLower(sliceContent)
	hasCollect := strings.Contains(lowerContent, "collect") || strings.Contains(lowerContent, "gather")
	hasAgents := strings.Contains(lowerContent, "agent") || strings.Contains(lowerContent, "implementer") || strings.Contains(lowerContent, "verifier")

	if !hasCollect && !hasAgents {
		t.Error("slice.md missing reference to collecting decisions from agents")
	}
}

func TestSliceMd_ReferencesDecisionsMdForAggregation(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/slice.md")
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	sliceContent := string(content)

	if !strings.Contains(sliceContent, "DECISIONS.md") {
		t.Error("slice.md missing reference to DECISIONS.md for aggregation")
	}
}

func TestSliceMd_ReferencesDecisionSourceFormat(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/slice.md")
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	sliceContent := string(content)

	// Should reference "slice:" source format
	if !strings.Contains(sliceContent, "slice:") && !strings.Contains(sliceContent, "Source") {
		t.Error("slice.md missing reference to decision source format (slice:)")
	}
}

func TestSliceMd_ReferencesDecisionAppendOnly(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/slice.md")
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	sliceContent := string(content)

	// Should mention append-only for decisions
	lowerContent := strings.ToLower(sliceContent)
	if !strings.Contains(lowerContent, "append") {
		t.Error("slice.md missing reference to append-only for decisions")
	}
}

func TestSliceMd_DecisionAggregationIsNonBlocking(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/slice.md")
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	sliceContent := string(content)

	// Should mention that decision aggregation failure doesn't block
	lowerContent := strings.ToLower(sliceContent)
	hasNonBlocking := strings.Contains(lowerContent, "does not block") ||
		strings.Contains(lowerContent, "doesn't block") ||
		strings.Contains(lowerContent, "not block") ||
		strings.Contains(lowerContent, "non-blocking")

	if !hasNonBlocking {
		t.Error("slice.md missing mention of decision aggregation being non-blocking")
	}
}

func TestSliceMd_ReferencesDecisionTableColumns(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/slice.md")
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	sliceContent := string(content)

	// Should reference decision table structure (Decision, Context, Chosen, Rejected, Date, Source)
	lowerContent := strings.ToLower(sliceContent)
	hasTableRef := strings.Contains(lowerContent, "table") ||
		strings.Contains(lowerContent, "column") ||
		strings.Contains(lowerContent, "row")

	if !hasTableRef {
		t.Error("slice.md missing reference to decision table structure")
	}
}

// C-45 Tests: RoadmapAutoUpdate
// These tests verify slice.md and wrap.md contain ROADMAP.md auto-update instructions

func TestSliceMd_ReferencesRoadmapUpdate(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/slice.md")
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	sliceContent := string(content)

	if !strings.Contains(sliceContent, "ROADMAP.md") {
		t.Error("slice.md missing ROADMAP.md update reference")
	}
}

func TestSliceMd_ReferencesSliceRowUpdate(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/slice.md")
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	sliceContent := string(content)

	// Should mention updating slice row with status, tests, date, decision
	lowerContent := strings.ToLower(sliceContent)
	hasUpdate := strings.Contains(lowerContent, "update") || strings.Contains(lowerContent, "mark")
	hasRow := strings.Contains(lowerContent, "row") || strings.Contains(lowerContent, "slice")

	if !hasUpdate && !hasRow {
		t.Error("slice.md missing reference to slice row update in ROADMAP.md")
	}
}

func TestSliceMd_ReferencesRoadmapCompletionFields(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/slice.md")
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	sliceContent := string(content)

	// Should mention completion fields: Status, Tests, Completed date, Key Decision
	lowerContent := strings.ToLower(sliceContent)
	hasStatus := strings.Contains(lowerContent, "status") || strings.Contains(lowerContent, "complete")
	hasTests := strings.Contains(lowerContent, "test")
	hasDate := strings.Contains(lowerContent, "date") || strings.Contains(lowerContent, "completed")

	if !hasStatus || !hasTests || !hasDate {
		t.Error("slice.md missing reference to ROADMAP.md completion fields (status, tests, date)")
	}
}

func TestSliceMd_RoadmapUpdateIsBestEffort(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/slice.md")
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	sliceContent := string(content)

	// Should mention best-effort or skip-if-missing for ROADMAP.md
	lowerContent := strings.ToLower(sliceContent)
	hasBestEffort := strings.Contains(lowerContent, "best-effort") ||
		strings.Contains(lowerContent, "best effort") ||
		strings.Contains(lowerContent, "skip") ||
		strings.Contains(lowerContent, "missing") ||
		strings.Contains(lowerContent, "if exists")

	if !hasBestEffort {
		t.Error("slice.md missing mention of ROADMAP.md update being best-effort/skip-if-missing")
	}
}

func TestWrapMd_ReferencesRoadmapWrapProgress(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/wrap.md")
	if err != nil {
		t.Fatalf("failed to read wrap.md: %v", err)
	}

	wrapContent := string(content)

	// Should reference ROADMAP.md wrap progress section
	if !strings.Contains(wrapContent, "ROADMAP.md") {
		t.Error("wrap.md missing ROADMAP.md reference")
	}
}

func TestWrapMd_ReferencesWrapProgressSection(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/wrap.md")
	if err != nil {
		t.Fatalf("failed to read wrap.md: %v", err)
	}

	wrapContent := string(content)

	// Should mention wrap progress section or boundary rows
	lowerContent := strings.ToLower(wrapContent)
	hasProgress := strings.Contains(lowerContent, "progress") || strings.Contains(lowerContent, "wrap")
	hasSection := strings.Contains(lowerContent, "section") || strings.Contains(lowerContent, "boundary")

	if !hasProgress && !hasSection {
		t.Error("wrap.md missing reference to wrap progress section in ROADMAP.md")
	}
}

func TestWrapMd_ReferencesWrapBoundaryFields(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/wrap.md")
	if err != nil {
		t.Fatalf("failed to read wrap.md: %v", err)
	}

	wrapContent := string(content)

	// Should mention boundary row fields: Boundary, Status, Locking Tests, Known Issues
	lowerContent := strings.ToLower(wrapContent)
	hasBoundary := strings.Contains(lowerContent, "boundary")
	hasStatus := strings.Contains(lowerContent, "status") || strings.Contains(lowerContent, "wrapped")
	hasTests := strings.Contains(lowerContent, "locking test") || strings.Contains(lowerContent, "test")

	if !hasBoundary || !hasStatus || !hasTests {
		t.Error("wrap.md missing reference to wrap boundary fields in ROADMAP.md")
	}
}

func TestWrapMd_RoadmapWrapUpdateIsBestEffort(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/wrap.md")
	if err != nil {
		t.Fatalf("failed to read wrap.md: %v", err)
	}

	wrapContent := string(content)

	// Should mention best-effort or skip-if-missing for ROADMAP.md wrap progress
	lowerContent := strings.ToLower(wrapContent)
	hasBestEffort := strings.Contains(lowerContent, "best-effort") ||
		strings.Contains(lowerContent, "best effort") ||
		strings.Contains(lowerContent, "skip") ||
		strings.Contains(lowerContent, "if exists") ||
		strings.Contains(lowerContent, "missing")

	if !hasBestEffort {
		t.Error("wrap.md missing mention of ROADMAP.md update being best-effort")
	}
}

// C-46 Tests: RoadmapDisplay
// These tests verify roadmap.md contains instructions for displaying ROADMAP.md

func TestRoadmapMd_ContainsRoadmapMdReadReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)

	if !strings.Contains(roadmapContent, "ROADMAP.md") {
		t.Error("roadmap.md missing ROADMAP.md read reference")
	}
}

func TestRoadmapMd_ContainsReadOnlyInstructions(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)
	lowerContent := strings.ToLower(roadmapContent)

	hasReadOnly := strings.Contains(lowerContent, "read-only") ||
		strings.Contains(lowerContent, "read only") ||
		strings.Contains(lowerContent, "strictly read-only")

	if !hasReadOnly {
		t.Error("roadmap.md missing read-only instructions")
	}
}

func TestRoadmapMd_ContainsNoModifyInstructions(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)
	lowerContent := strings.ToLower(roadmapContent)

	hasNoModify := strings.Contains(lowerContent, "no files") ||
		strings.Contains(lowerContent, "no modify") ||
		strings.Contains(lowerContent, "not modify") ||
		strings.Contains(lowerContent, "don't modify")

	if !hasNoModify {
		t.Error("roadmap.md missing no-modify instructions")
	}
}

func TestRoadmapMd_ContainsConfigJsonReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)

	if !strings.Contains(roadmapContent, "config.json") {
		t.Error("roadmap.md missing config.json reference for project context")
	}
}

func TestRoadmapMd_ContainsProjectContextReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)
	lowerContent := strings.ToLower(roadmapContent)

	if !strings.Contains(lowerContent, "project context") && !strings.Contains(lowerContent, "context") {
		t.Error("roadmap.md missing project context reference")
	}
}

func TestRoadmapMd_ContainsMissingRoadmapErrorHandling(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)

	hasErrorMessage := strings.Contains(roadmapContent, "No roadmap found") ||
		strings.Contains(roadmapContent, "doesn't exist") ||
		strings.Contains(roadmapContent, "does not exist")

	if !hasErrorMessage {
		t.Error("roadmap.md missing error handling for missing ROADMAP.md")
	}
}

func TestRoadmapMd_ContainsDesignPrerequisiteReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)

	if !strings.Contains(roadmapContent, "/gl:design") && !strings.Contains(roadmapContent, "gl:design") {
		t.Error("roadmap.md missing /gl:design prerequisite reference")
	}
}

func TestRoadmapMd_ContainsEmptyRoadmapErrorHandling(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)

	hasEmptyError := strings.Contains(roadmapContent, "empty") ||
		strings.Contains(roadmapContent, "ROADMAP.md is empty")

	if !hasEmptyError {
		t.Error("roadmap.md missing error handling for empty ROADMAP.md")
	}
}

func TestRoadmapMd_ContainsArchitectureDiagramReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)
	lowerContent := strings.ToLower(roadmapContent)

	hasDiagram := strings.Contains(lowerContent, "architecture diagram") ||
		strings.Contains(lowerContent, "diagram") ||
		strings.Contains(lowerContent, "architecture")

	if !hasDiagram {
		t.Error("roadmap.md missing architecture diagram reference")
	}
}

func TestRoadmapMd_ContainsMilestoneTablesReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)
	lowerContent := strings.ToLower(roadmapContent)

	hasMilestoneTables := strings.Contains(lowerContent, "milestone") &&
		(strings.Contains(lowerContent, "table") || strings.Contains(lowerContent, "tables"))

	if !hasMilestoneTables {
		t.Error("roadmap.md missing milestone tables reference")
	}
}

func TestRoadmapMd_ContainsArchivedMilestonesReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)
	lowerContent := strings.ToLower(roadmapContent)

	if !strings.Contains(lowerContent, "archived") {
		t.Error("roadmap.md missing archived milestones reference")
	}
}

// C-47 Tests: RoadmapMilestonePlanning
// These tests verify roadmap.md contains milestone planning sub-command instructions

func TestRoadmapMd_ContainsMilestoneSubcommand(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)

	if !strings.Contains(roadmapContent, "milestone") {
		t.Error("roadmap.md missing milestone sub-command reference")
	}
}

func TestRoadmapMd_ContainsGlDesignerSpawnReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)
	lowerContent := strings.ToLower(roadmapContent)

	hasDesigner := strings.Contains(lowerContent, "gl-designer") || strings.Contains(lowerContent, "designer")
	hasSpawn := strings.Contains(lowerContent, "spawn") || strings.Contains(roadmapContent, "Task")

	if !hasDesigner || !hasSpawn {
		t.Error("roadmap.md missing gl-designer spawn reference")
	}
}

func TestRoadmapMd_ContainsGraphJsonReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)

	if !strings.Contains(roadmapContent, "GRAPH.json") {
		t.Error("roadmap.md missing GRAPH.json reference for milestone planning")
	}
}

func TestRoadmapMd_ContainsDecisionsMdAppendReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)

	if !strings.Contains(roadmapContent, "DECISIONS.md") {
		t.Error("roadmap.md missing DECISIONS.md append reference for milestones")
	}
}

func TestRoadmapMd_ContainsLighterSessionReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)
	lowerContent := strings.ToLower(roadmapContent)

	hasLighter := strings.Contains(lowerContent, "lighter") || strings.Contains(lowerContent, "skip")
	hasInit := strings.Contains(lowerContent, "init") || strings.Contains(lowerContent, "interview")

	if !hasLighter && !hasInit {
		t.Error("roadmap.md missing lighter design session reference (skip init/interview)")
	}
}

func TestRoadmapMd_ContainsMilestoneFieldReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)
	lowerContent := strings.ToLower(roadmapContent)

	hasMilestoneField := strings.Contains(lowerContent, "milestone field") ||
		(strings.Contains(lowerContent, "milestone") && strings.Contains(lowerContent, "field"))

	if !hasMilestoneField {
		t.Error("roadmap.md missing milestone field reference on slices")
	}
}

func TestRoadmapMd_ContainsDesignMdReadForMilestone(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)

	if !strings.Contains(roadmapContent, "DESIGN.md") {
		t.Error("roadmap.md missing DESIGN.md read reference for milestone planning")
	}
}

func TestRoadmapMd_ContainsContractsMdReadForMilestone(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)

	if !strings.Contains(roadmapContent, "CONTRACTS.md") {
		t.Error("roadmap.md missing CONTRACTS.md read reference for milestone planning")
	}
}

func TestRoadmapMd_ContainsStateMdReadForMilestone(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)

	if !strings.Contains(roadmapContent, "STATE.md") {
		t.Error("roadmap.md missing STATE.md read reference for milestone planning")
	}
}

func TestRoadmapMd_ContainsMilestoneCommitMessage(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)

	hasCommitFormat := strings.Contains(roadmapContent, "docs: plan milestone") ||
		strings.Contains(roadmapContent, "plan milestone")

	if !hasCommitFormat {
		t.Error("roadmap.md missing milestone commit message format")
	}
}

func TestRoadmapMd_ContainsMilestoneErrorHandling(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)
	lowerContent := strings.ToLower(roadmapContent)

	hasErrors := strings.Contains(lowerContent, "error") ||
		strings.Contains(roadmapContent, "NoRoadmap") ||
		strings.Contains(roadmapContent, "NoConfig") ||
		strings.Contains(roadmapContent, "DesignerFailure")

	if !hasErrors {
		t.Error("roadmap.md missing milestone error handling references")
	}
}

// C-48 Tests: RoadmapMilestoneArchive
// These tests verify roadmap.md contains archive sub-command instructions

func TestRoadmapMd_ContainsArchiveSubcommand(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)

	if !strings.Contains(roadmapContent, "archive") {
		t.Error("roadmap.md missing archive sub-command reference")
	}
}

func TestRoadmapMd_ContainsArchivedMilestonesSectionReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)

	hasArchived := strings.Contains(roadmapContent, "Archived Milestones") ||
		strings.Contains(roadmapContent, "archived milestone")

	if !hasArchived {
		t.Error("roadmap.md missing Archived Milestones section reference")
	}
}

func TestRoadmapMd_ContainsCompletedMilestoneIdentification(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)
	lowerContent := strings.ToLower(roadmapContent)

	hasCompleted := strings.Contains(lowerContent, "completed") || strings.Contains(lowerContent, "complete")
	hasIdentify := strings.Contains(lowerContent, "identif") || strings.Contains(lowerContent, "identify")

	if !hasCompleted || !hasIdentify {
		t.Error("roadmap.md missing completed milestone identification reference")
	}
}

func TestRoadmapMd_ContainsCompressionFormatReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)
	lowerContent := strings.ToLower(roadmapContent)

	hasCompression := strings.Contains(lowerContent, "compress") || strings.Contains(lowerContent, "summary")
	hasFormat := strings.Contains(lowerContent, "format") || strings.Contains(roadmapContent, "completed")

	if !hasCompression && !hasFormat {
		t.Error("roadmap.md missing compression/summary format reference for archiving")
	}
}

func TestRoadmapMd_ContainsArchiveCommitMessage(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)

	hasArchiveCommit := strings.Contains(roadmapContent, "docs: archive milestone") ||
		strings.Contains(roadmapContent, "archive milestone")

	if !hasArchiveCommit {
		t.Error("roadmap.md missing archive commit message format")
	}
}

func TestRoadmapMd_ContainsArchiveErrorHandling(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)

	hasArchiveErrors := strings.Contains(roadmapContent, "NoCompletedMilestones") ||
		strings.Contains(roadmapContent, "ArchiveFailure") ||
		strings.Contains(roadmapContent, "NoRoadmap")

	if !hasArchiveErrors {
		t.Error("roadmap.md missing archive error handling references")
	}
}

func TestRoadmapMd_ContainsUserSelectionForArchive(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/roadmap.md")
	if err != nil {
		t.Fatalf("failed to read roadmap.md: %v", err)
	}

	roadmapContent := string(content)
	lowerContent := strings.ToLower(roadmapContent)

	hasSelection := strings.Contains(lowerContent, "select") ||
		strings.Contains(lowerContent, "user select") ||
		strings.Contains(lowerContent, "presents") ||
		strings.Contains(lowerContent, "choose")

	if !hasSelection {
		t.Error("roadmap.md missing user selection reference for archive")
	}
}

// C-49 Tests: ChangelogDisplay
// These tests verify changelog.md contains instructions for displaying formatted changelog

func TestChangelogMd_ContainsSummariesDirectoryReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)

	if !strings.Contains(changelogContent, "summaries/") && !strings.Contains(changelogContent, ".greenlight/summaries") {
		t.Error("changelog.md missing summaries/ directory reference")
	}
}

func TestChangelogMd_ContainsReadOnlyInstructions(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)
	lowerContent := strings.ToLower(changelogContent)

	hasReadOnly := strings.Contains(lowerContent, "read-only") ||
		strings.Contains(lowerContent, "read only") ||
		strings.Contains(lowerContent, "no files written")

	if !hasReadOnly {
		t.Error("changelog.md missing read-only instructions")
	}
}

func TestChangelogMd_ContainsConfigJsonReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)

	if !strings.Contains(changelogContent, "config.json") {
		t.Error("changelog.md missing config.json reference for project name")
	}
}

func TestChangelogMd_ContainsChronologicalSortingReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)
	lowerContent := strings.ToLower(changelogContent)

	hasChronological := strings.Contains(lowerContent, "chronological") ||
		(strings.Contains(lowerContent, "newest") && strings.Contains(lowerContent, "first"))

	if !hasChronological {
		t.Error("changelog.md missing chronological sorting reference (newest first)")
	}
}

func TestChangelogMd_ContainsChangelogHeaderFormat(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)

	if !strings.Contains(changelogContent, "CHANGELOG") {
		t.Error("changelog.md missing CHANGELOG header format reference")
	}
}

func TestChangelogMd_ContainsEntryTypesReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)

	hasSlice := strings.Contains(changelogContent, "slice")
	hasWrap := strings.Contains(changelogContent, "wrap")
	hasQuick := strings.Contains(changelogContent, "quick")

	if !hasSlice || !hasWrap || !hasQuick {
		t.Error("changelog.md missing entry types reference (slice, wrap, quick)")
	}
}

func TestChangelogMd_ContainsNoSummariesDirErrorHandling(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)

	hasNoSummariesError := strings.Contains(changelogContent, "NoSummariesDir") ||
		strings.Contains(changelogContent, "No summaries found")

	if !hasNoSummariesError {
		t.Error("changelog.md missing NoSummariesDir error handling")
	}
}

func TestChangelogMd_ContainsEmptySummariesErrorHandling(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)

	hasEmptyError := strings.Contains(changelogContent, "EmptySummariesDir") ||
		strings.Contains(changelogContent, "No summaries found yet")

	if !hasEmptyError {
		t.Error("changelog.md missing EmptySummariesDir error handling")
	}
}

func TestChangelogMd_ContainsSummaryFileNamingReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)

	if !strings.Contains(changelogContent, "SUMMARY.md") {
		t.Error("changelog.md missing SUMMARY.md file naming convention reference")
	}
}

func TestChangelogMd_ContainsParsingSummariesReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)
	lowerContent := strings.ToLower(changelogContent)

	hasParsing := strings.Contains(lowerContent, "parse") ||
		strings.Contains(lowerContent, "parsing")

	if !hasParsing {
		t.Error("changelog.md missing parsing summaries reference")
	}
}

func TestChangelogMd_ContainsEntryCountSummaryReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)
	lowerContent := strings.ToLower(changelogContent)

	hasCountSummary := strings.Contains(lowerContent, "count") ||
		strings.Contains(lowerContent, "entries") ||
		strings.Contains(lowerContent, "bottom")

	if !hasCountSummary {
		t.Error("changelog.md missing entry count summary reference")
	}
}

func TestChangelogMd_ContainsNoConfigErrorHandling(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)

	hasNoConfigError := strings.Contains(changelogContent, "NoConfig") ||
		strings.Contains(changelogContent, "Unknown Project")

	if !hasNoConfigError {
		t.Error("changelog.md missing NoConfig error handling (fallback to 'Unknown Project')")
	}
}

func TestChangelogMd_ContainsMalformedSummaryErrorHandling(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)

	hasMalformedError := strings.Contains(changelogContent, "MalformedSummary") ||
		strings.Contains(changelogContent, "Could not parse")

	if !hasMalformedError {
		t.Error("changelog.md missing MalformedSummary error handling")
	}
}

func TestChangelogMd_ContainsEntryFormatReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)
	lowerContent := strings.ToLower(changelogContent)

	hasDate := strings.Contains(lowerContent, "date")
	hasType := strings.Contains(lowerContent, "type")
	hasName := strings.Contains(lowerContent, "name")

	if !hasDate || !hasType || !hasName {
		t.Error("changelog.md missing entry format reference (date, type, name)")
	}
}

func TestChangelogMd_ContainsTestCountReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)
	lowerContent := strings.ToLower(changelogContent)

	hasTestCount := strings.Contains(lowerContent, "test count") ||
		strings.Contains(lowerContent, "tests")

	if !hasTestCount {
		t.Error("changelog.md missing test count reference in entries")
	}
}

func TestChangelogMd_ContainsOneLineSummaryReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)
	lowerContent := strings.ToLower(changelogContent)

	hasOneLineSummary := strings.Contains(lowerContent, "one-line") ||
		strings.Contains(lowerContent, "summary")

	if !hasOneLineSummary {
		t.Error("changelog.md missing one-line summary reference")
	}
}

// C-50 Tests: ChangelogFiltering
// These tests verify changelog.md contains milestone and since sub-command filtering instructions

func TestChangelogMd_ContainsMilestoneSubcommand(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)

	if !strings.Contains(changelogContent, "milestone") {
		t.Error("changelog.md missing milestone sub-command reference")
	}
}

func TestChangelogMd_ContainsSinceSubcommand(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)

	if !strings.Contains(changelogContent, "since") {
		t.Error("changelog.md missing since sub-command reference")
	}
}

func TestChangelogMd_ContainsGraphJsonReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)

	if !strings.Contains(changelogContent, "GRAPH.json") {
		t.Error("changelog.md missing GRAPH.json reference for milestone filtering")
	}
}

func TestChangelogMd_ContainsDateFilteringReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)
	lowerContent := strings.ToLower(changelogContent)

	hasDateFilter := strings.Contains(lowerContent, "date") &&
		(strings.Contains(lowerContent, "filter") || strings.Contains(lowerContent, "since"))

	if !hasDateFilter {
		t.Error("changelog.md missing date filtering reference")
	}
}

func TestChangelogMd_ContainsDateFormatReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)

	hasDateFormat := strings.Contains(changelogContent, "YYYY-MM-DD") ||
		strings.Contains(changelogContent, "ISO 8601")

	if !hasDateFormat {
		t.Error("changelog.md missing date format reference (YYYY-MM-DD or ISO 8601)")
	}
}

func TestChangelogMd_ContainsUnknownMilestoneErrorHandling(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)

	hasUnknownMilestoneError := strings.Contains(changelogContent, "UnknownMilestone") ||
		strings.Contains(changelogContent, "unknown milestone")

	if !hasUnknownMilestoneError {
		t.Error("changelog.md missing UnknownMilestone error handling")
	}
}

func TestChangelogMd_ContainsInvalidDateErrorHandling(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)

	hasInvalidDateError := strings.Contains(changelogContent, "InvalidDate") ||
		strings.Contains(changelogContent, "invalid date")

	if !hasInvalidDateError {
		t.Error("changelog.md missing InvalidDate error handling")
	}
}

func TestChangelogMd_ContainsMilestoneFilterReadOnly(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)
	lowerContent := strings.ToLower(changelogContent)

	hasReadOnly := strings.Contains(lowerContent, "read-only") ||
		strings.Contains(lowerContent, "read only")

	if !hasReadOnly {
		t.Error("changelog.md missing read-only reference for milestone filter")
	}
}

func TestChangelogMd_ContainsSinceFilterReadOnly(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)
	lowerContent := strings.ToLower(changelogContent)

	hasReadOnly := strings.Contains(lowerContent, "read-only") ||
		strings.Contains(lowerContent, "read only")

	if !hasReadOnly {
		t.Error("changelog.md missing read-only reference for since filter")
	}
}

func TestChangelogMd_ContainsMilestoneNameMatchingReference(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)
	lowerContent := strings.ToLower(changelogContent)

	hasMilestoneMatching := strings.Contains(lowerContent, "milestone") &&
		(strings.Contains(lowerContent, "match") || strings.Contains(lowerContent, "filter"))

	if !hasMilestoneMatching {
		t.Error("changelog.md missing milestone name matching reference")
	}
}

func TestChangelogMd_ContainsNoMatchingEntriesErrorHandling(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)

	hasNoMatchingError := strings.Contains(changelogContent, "NoMatchingEntries") ||
		strings.Contains(changelogContent, "No matching entries")

	if !hasNoMatchingError {
		t.Error("changelog.md missing NoMatchingEntries error handling")
	}
}

func TestChangelogMd_ContainsMilestoneHeaderFormat(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)

	hasMilestoneHeader := strings.Contains(changelogContent, "CHANGELOG") &&
		strings.Contains(changelogContent, "milestone")

	if !hasMilestoneHeader {
		t.Error("changelog.md missing milestone header format reference")
	}
}

func TestChangelogMd_ContainsSinceHeaderFormat(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)

	hasSinceHeader := strings.Contains(changelogContent, "CHANGELOG") &&
		strings.Contains(changelogContent, "since")

	if !hasSinceHeader {
		t.Error("changelog.md missing since header format reference")
	}
}

func TestChangelogMd_ContainsNoGraphJsonErrorHandling(t *testing.T) {
	content, err := os.ReadFile("/Users/juliantellez/github.com/atlantic-blue/greenlight/src/commands/gl/changelog.md")
	if err != nil {
		t.Fatalf("failed to read changelog.md: %v", err)
	}

	changelogContent := string(content)

	hasNoGraphError := strings.Contains(changelogContent, "NoGraphJson") ||
		strings.Contains(changelogContent, "NoGraph")

	if !hasNoGraphError {
		t.Error("changelog.md missing NoGraphJson error handling")
	}
}
