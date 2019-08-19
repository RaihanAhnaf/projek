var metadataVM = kendo.observable({
	arraj: [],
	url: {
		update: "/upload/updatemetadata"
	},
	allConfirmed: false,
	allConfirmOrTrained: false,
	showNextStep: false,
	init: function(){
		var self = this

		var rows = self.get("arraj")
		_.each(rows, function(row){
			if(reviewVM.get("statussummary") == "New Contract" && (row.get("contentFromUser") == null || row.get("contentFromUser") == "Invalid Date" || row.get("contentFromUser") == "")){
				if(row.get("type") == "date"){
					row.set("contentFromUser", getDateFromString(row.get("contentFromSystem")).date)
				} else if(row.get("type") == "dropdown"){
					row.set("contentFromUser", getFuzzyResult(row.get("contentFromSystem")))
				} else {
					row.set("contentFromUser", row.get("contentFromSystem"))
				}
			}

		    reviewVM.set("updateParam", [{
		        "confidence": 1.0,
		        "clause": row.get("clauseid"),
		        "contentfromuser": row.get("contentFromUser"),
		        "contentfromsystem": row.get("contentFromSystem"),
		        "common-ancestor-id": row.get("metadata").length > 0 ? row.get("metadata")[row.get("metadata").length - 1]["commonancestorid"] : "",
		        "data-eaciit-id": row.get("metadata").length > 0 ? row.get("metadata")[row.get("metadata").length - 1]["dataeaciitid"] : "",
		        "start-offset": row.get("metadata").length > 0 ? row.get("metadata")[row.get("metadata").length - 1]["startoffset"] : "",
		        "end-offset": row.get("metadata").length > 0 ? row.get("metadata")[row.get("metadata").length - 1]["endoffset"] : "",
		        "data-eaciit-start": row.get("metadata").length > 0 ? row.get("metadata")[row.get("metadata").length - 1]["dataeaciitstart"] : "",
		        "data-eaciit-end": row.get("metadata").length > 0 ? row.get("metadata")[row.get("metadata").length - 1]["dataeaciitend"] : "",
		        "identified": row.get("isIdentified"),
		        "confirmortrained": row.get("confirmOrTrained")
		    }])

		    ajaxPost("/upload/updatemetadata", { 
		        id: reviewVM.get("id"),
		        clause: row.get("clauseid"),
		        data: reviewVM.get("updateParam")
		    }, function(resUpdate){  })

			self.setConfirmStatus(row)
		});

		leftCol.checkAllConfirmed()

		var subcluster = reviewVM.get("subcluster")
		self.set("showNextStep", (subcluster == "externalmsa" || subcluster == "externalsla"))
	},
	setThisRowToActive: function(e){
		leftCol.setRowToActive(e.data)
	},
	setConfirmStatus: function(row){
		row.set("confirmStatus", (function(){
			if((row.get("isIdentified"))){
				return row.get("confirmOrTrained") == "Trained" ? "Trained" : (row.get("confirmOrTrained") == "Not Found" ? "Unidentified" : "Identified")
			} else {
				if(
					// reviewVM.get("statussummary") == "New Contract" && 
					((row.get("contentFromSystem") != null && row.get("contentFromSystem") != "Invalid Date" && row.get("contentFromSystem") != "" && row.get("contentFromSystem") != undefined))){
					return 'Identified';
				} else {
					return 'Unidentified';
				}
			}
		})())
	}
})

metadataVM.hardConfirmVisible = function(row){
	if(row.get("confirmOrTrained") == "" && row.get("confirmStatus") == "Identified"){
		return true
	} else {
		return false
	}
}

metadataVM.hardConfirmClicked = function(e){
	var self = this

	var currentVM = leftCol.get("currentVM")

	var row = e.data

	row.set("isIdentified", true)
	row.set("confirmOrTrained", "Confirmed")

    reviewVM.set("updateParam", [{
        "confidence": 1.0,
        "clause": row.get("clauseid"),
        "contentfromuser": row.get("contentFromUser"),
        "contentfromsystem": row.get("contentFromSystem"),
        "common-ancestor-id": row.get("metadata").length > 0 ? row.get("metadata")[row.get("metadata").length - 1]["commonancestorid"] : "",
        "data-eaciit-id": row.get("metadata").length > 0 ? row.get("metadata")[row.get("metadata").length - 1]["dataeaciitid"] : "",
        "start-offset": row.get("metadata").length > 0 ? row.get("metadata")[row.get("metadata").length - 1]["startoffset"] : "",
        "end-offset": row.get("metadata").length > 0 ? row.get("metadata")[row.get("metadata").length - 1]["endoffset"] : "",
        "data-eaciit-start": row.get("metadata").length > 0 ? row.get("metadata")[row.get("metadata").length - 1]["dataeaciitstart"] : "",
        "data-eaciit-end": row.get("metadata").length > 0 ? row.get("metadata")[row.get("metadata").length - 1]["dataeaciitend"] : "",
        "identified": row.get("isIdentified"),
        "confirmortrained": row.get("confirmOrTrained")
    }])

    ajaxPost("/upload/updatemetadata", { 
        id: reviewVM.get("id"),
        clause: row.get("clauseid"),
        data: reviewVM.get("updateParam")
    }, function(resUpdate){
		leftCol.checkAllConfirmed()
		// self.setNextRowToActive()

		currentVM.setConfirmStatus(row)

		unSelect()
		// reviewVM.refresh()
		leftCol.setNextRowToActive()
    })
}

metadataVM.train = function(e){
	leftCol.trainRow(e.data)
}

metadataVM.setNF = function(e){ 
	leftCol.setRowNotFound(e.data)
}

metadataVM.inputChange = function(e){
	if(leftCol.get("currentVM").allConfirmed == true)
		e.data.set("changesAfterAllIsConfirmed", true)
}

metadataVM.isDropdown = function(row) {
	return (row.get("type") == "dropdown") ? true : false
}

metadataVM.isDate = function(row) {
	return (row.get("type") == "date") ? true : false
}

metadataVM.isText = function(row) {
	return (row.get("type") == null) ? true : false
}

metadataVM.isConfirmLabel = function(row){
	if(row.get("isIdentified")){
		return true
	} else{
		return false
	}
}

metadataVM.confirmLabelText = function(row){
	if(row.get("isIdentified")){
		if(row.get("confirmOrTrained") == "Not Found")
			return "Not Found"
		else
			return "Confirmed"
	} else{
		return ""
	}
}

metadataVM.actionAfterActionVisibility = function(row) {
	if(row.get("isIdentified") 
		&& row.get("confirmOrTrained") == "Trained" 
		&& ((reviewVM.get("statussummary") == "New Contract" && row.get("confirmStatus") == "Unidentified" && row.get("isConfirmed") == false) 
			|| (reviewVM.get("statussummary") == "New Contract" && row.get("confirmStatus") == "Identified" && row.get("isConfirmed") == false)
			|| (reviewVM.get("statussummary") != "New Contract" && row.get("confirmStatus") == "Unidentified" && row.get("isConfirmed") == false)
			|| (reviewVM.get("statussummary") != "New Contract" && row.get("confirmStatus") == "Identified" && row.get("isConfirmed") == false)
			)
		){
		return true
	} else {
		return false
	}
}

metadataVM.cancelRow = function(e){
	e.data.set("confirmOrTrained", "")
	e.data.set("isIdentified", false)
	e.data.set("contentFromSystem", "")
	e.data.set("contentFromUser", "")
	e.data.set("metadata", [])

    viewerVM.set("highlightClause", "")

	leftCol.checkAllConfirmed();
}

metadataVM.confirmRow = function(e){
	var self = this

	var updateParam = reviewVM.get("updateParam")
	var currentRow = leftCol.get("currentRow")

	_.each(updateParam, function(up){
		if (up.get("contentfromuser") == null || up.get("contentfromuser") == "" || up.get("contentfromuser") == undefined || up.get("contentfromuser") != currentRow.get("contentFromUser"))
			up.set("contentfromuser", currentRow.get("contentFromUser"))
		if ((up.get("contentfromuserdesc") == null || up.get("contentfromuserdesc") == "" || up.get("contentfromuserdesc") == undefined) && currentRow.get("type") == 'dropdown')
			up.set("contentfromuserdesc", currentRow.get("contentFromUser"))
	})

	ajaxPost("/upload/updatemetadata", { 
        id: reviewVM.get("id"),
        clause: currentRow.get("clauseid"),
        data: updateParam
    }, function(resUpdate){
        reviewVM.refresh(leftCol.getNextRowIndex())
    })
}

metadataVM.nextStep = function(e){
	var status = "META"
	ajaxPost("/upload/updatemetadata", { 
        id: reviewVM.get("id"),
        clause: e.data.get("clauseid"),
        data: reviewVM.get("updateParam"),
        status: status
    }, function(resUpdate){
    	reviewVM.refresh()
    })
}

metadataVM.saveNClose = function(e){
	var self = this

	var metadataUpdateAfterConfirmed = {
		id: reviewVM.get("id"),
		data: _.map(_.filter(self.get("arraj"), function(v){ 
			return v.changesAfterAllIsConfirmed == true 
		}), function(a){
			_.each(a.metadata, function(meta){
				meta.contentfromuser = a.contentFromUser
			})
			return {
				clause: a.clauseid,
				metadata: a.metadata
			}
		})
	}

	var url = "", param = {}
	if(metadataUpdateAfterConfirmed.data.length > 0){
		url = "/upload/updatecontentfromuser"
		param = metadataUpdateAfterConfirmed
	} else {
		url = "/upload/updatemetadata"
		param = { 
	        id: reviewVM.get("id"),
	        clause: e.data.clauseid,
	        data: reviewVM.get("updateParam"),
	        status: leftCol.areWeTakingTheQuiz() ? "META" : ""
	    }
	}

	ajaxPost(url, param, function(resUpdate){
	    // leftCol.checkAllConfirmed();
	    // reviewVM.refresh()

        window.location = "/dashboard/default?p=2";
    })
}

metadataVM.captureHighlighted = function(){
	var self = this

	var highlightClause = viewerVM.get("highlightClause")

    if(highlightClause != ""){

        ((function() {
            var selected = doGetSelection()

            if(viewerVM.highlightValidation(selected) == true){
                var currentRow = leftCol.get("currentRow")

                var selectedText = captureSelection()

                currentRow.set("contentFromUser", (self.isText(currentRow)) ? selectedText : currentRow.get("contentFromUser"))
                currentRow.set("contentFromSystem", selectedText)
                currentRow.set("isIdentified", true)
                currentRow.set("confirmOrTrained", "Trained")

                if (currentRow.get("clauseid") == "governinglaw"){
                    var res = getFuzzyResult(currentRow.get("contentFromSystem"))
                    
                    if (res != null && res != ""){
                        currentRow.set("contentFromUser", res)
                    } else {
                        swal("Warning!", "Cannot detect country name from selection.", "warning");
                    }

                    reviewVM.set("updateParam", [{
                        "confidence": 1.0,
                        "clause": highlightClause,
                        "contentfromuser": currentRow.get("contentFromUser"),
                        "contentfromuserdesc": currentRow.get("contentFromUser"),
                        "contentfromsystem": currentRow.get("contentFromSystem"),
                        "common-ancestor-id": selected.commonAncestorId,
                        "data-eaciit-id": "",
                        "data-eaciit-start": selected.startElementId,
                        "data-eaciit-end": selected.endElementId,
                        "start-offset": selected.startOffset.toString(),
                        "end-offset": selected.endOffset.toString(),
                        "identified": currentRow.get("isIdentified"),
                        "confirmortrained": currentRow.get("confirmOrTrained")
                    }])

                } else if(currentRow.get("clauseid") == "effectivedate"){
                    highlightstr = currentRow.get("contentFromSystem")
                    datestr = getDateFromString(highlightstr)

                    if(datestr != undefined){
	                    if(datestr.date != null){
							currentRow.set("contentFromUser", datestr.date)
	                    }else{
							swal("Warning!", "Selection Not Contains Date", 'warning');
						}
                    } else {
                    	swal("Warning!", "Cannot detect date from selection", 'warning');
                    }

                    reviewVM.set("updateParam", [{
                        "confidence": 1.0,
                        "clause": highlightClause,
                        "contentfromuser": currentRow.get("contentFromUser"),
                        "contentfromsystem": currentRow.get("contentFromSystem"),
                        "common-ancestor-id": selected.commonAncestorId,
                        "data-eaciit-id": "",
                        "data-eaciit-start": selected.startElementId,
                        "data-eaciit-end": selected.endElementId,
                        "start-offset": selected.startOffset.toString(),
                        "end-offset": selected.endOffset.toString(),
                        "identified": currentRow.get("isIdentified"),
                        "confirmortrained": currentRow.get("confirmOrTrained")
                    }])
                } else {
                    reviewVM.set("updateParam", [{
                        "confidence": 1.0,
                        "clause": highlightClause,
                        "contentfromuser": currentRow.get("contentFromUser"),
                        "contentfromsystem": currentRow.get("contentFromSystem"),
                        "common-ancestor-id": selected.commonAncestorId,
                        "data-eaciit-id": "",
                        "data-eaciit-start": selected.startElementId,
                        "data-eaciit-end": selected.endElementId,
                        "start-offset": selected.startOffset.toString(),
                        "end-offset": selected.endOffset.toString(),
                        "identified": currentRow.get("isIdentified"),
                        "confirmortrained": currentRow.get("confirmOrTrained")
                    }])
                }

                currentRow.set("metadata", reviewVM.get("updateParam"))

                viewerVM.set("highlightClause", "")
            }
        })())
    }
}

$(document).ready(function(){
	kendo.bind($("#metadataPage"), metadataVM);

	setTimeout(function(){
		$("#metadataPage .btn.information").tooltipster({
			trigger: 'hover'
		});
	}, 1000)
})