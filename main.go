package main

import (
	"fmt"
	"log"
	"os"

	"robot-readme/openapi" // Replace with your actual module name if different
)

func main() {
	specPath := "swagger.json"

	log.Printf("Reading spec from: %s\n", specPath)
	doc, err := openapi.LoadAPISpec(specPath)
	if err != nil {
		log.Fatalf("Error loading API spec: %v", err)
	}

	// Quick debug: print some top-level info from doc
	log.Printf("Loaded doc: Title=%s, Version=%s, #Endpoints=%d",
		doc.Title,
		doc.Version,
		len(doc.Endpoints),
	)

	if err := openapi.ResolveReferences(doc); err != nil {
		log.Fatalf("Error resolving references: %v", err)
	}

	summary := openapi.RenderText(doc)

	// Write to file
	outputFile := "llm1.txt"
	log.Printf("Writing summary to %s...", outputFile)

	err = os.WriteFile(outputFile, []byte(summary), 0644)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	fmt.Printf("Successfully wrote API summary to %s\n", outputFile)
}
