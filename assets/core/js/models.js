var metadataRows = { 
	documenttitle: { name: "Document Title", type: null },
	effectivedate: { name: "Effective Date", type: "date" },
	governinglaw: { name: "Governing Law", type: "dropdown", list: [] },
}

var questionnaireRows = {
	antibriberyclause: { name: "Anti-Bribery & Corruption", sentence: "Is there an ABC clause?" },
	protectionclause: { name: "Privacy", sentence: "Is there a data protection clause?" },
}

var metadataMapping = {
	internalsla: metadataRows,
	internalmsa: metadataRows,
	externalsla: metadataRows,
	externalmsa: metadataRows
}

var questionnaireMapping = {
	externalsla: questionnaireRows,
	externalmsa: questionnaireRows
}