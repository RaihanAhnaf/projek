var quizVM = kendo.observable({
	arraj: [],
	url: {
		update: "/upload/updatequiz"
	},
	allConfirmed: false,
	init: function(){
		var self = this

		var rows = self.get("arraj")
		_.each(rows, function(row){
			self.setConfirmStatus(row)
		});

		leftCol.checkAllConfirmed()
	},
	setThisRowToActive: function(e){
		leftCol.setRowToActive(e.data)
	},
	setConfirmStatus: function(row){
		row.set("confirmStatus", (function(){
			if((row.get("isIdentified"))){
				return row.get("confirmOrTrained") == "Trained" ? "Trained" : (row.get("confirmOrTrained") == "Not Found" ? "Unidentified" : "Identified")
			} else{
				return 'Unidentified';
			}
		})())
	}
})

quizVM.isConfirmLabel = function(row){
	if(row.get("isIdentified")){
		return true
	} else{
		return false
	}
}

quizVM.isYes = function(row){
    return row.get("contentFromUser") == "Yes"
}

quizVM.isNo = function(row){
    return row.get("contentFromUser") == "No"
}

quizVM.confirmLabelText = function(row){
	if(row.get("isIdentified")){
		if(row.get("confirmOrTrained") == "Not Found")
			return "Not Found"
		else
			return "Confirmed"
	} else{
		return ""
	}
}

quizVM.train = function(e){
	leftCol.trainRow(e.data)
}

quizVM.setNF = function(e){ 
	leftCol.setRowNotFound(e.data)
}

quizVM.yes = function(e){
	var self = this

	var row = e.data
	row.set("contentFromUser", "Yes")

    leftCol.trainRow(row)
}

quizVM.no = function(e){
	var self = this

	var row = e.data
	row.set("contentFromUser", "No")
	row.set("isIdentified", true)
	row.set("confirmOrTrained", "Confirmed")

    reviewVM.set("updateParam", [{
        "confidence": 1.0,
        "clause": row.get("clauseid"),
        "contentfromuser": row.get("contentFromUser"),
        "contentfromsystem": row.get("contentFromSystem"),
        "common-ancestor-id": row.get("metadata").length > 0 ? row.get("metadata")[row.get("metadata").length - 1]["common-ancestor-id"] : "",
        "data-eaciit-id": row.get("metadata").length > 0 ? row.get("metadata")[row.get("metadata").length - 1]["data-eaciit-id"] : "",
        "start-offset": row.get("metadata").length > 0 ? row.get("metadata")[row.get("metadata").length - 1]["start-offset"] : "",
        "end-offset": row.get("metadata").length > 0 ? row.get("metadata")[row.get("metadata").length - 1]["common-ancestor-id"] : "",
        "data-eaciit-start": row.get("metadata").length > 0 ? row.get("metadata")[row.get("metadata").length - 1]["data-eaciit-start"] : "",
        "data-eaciit-end": row.get("metadata").length > 0 ? row.get("metadata")[row.get("metadata").length - 1]["data-eaciit-end"] : "",
        "start-offset": row.get("metadata").length > 0 ? row.get("metadata")[row.get("metadata").length - 1]["start-offset"] : "",
        "identified": row.get("isIdentified"),
        "confirmortrained": row.get("confirmOrTrained")
    }])

    ajaxPost(self.url.update, { 
        id: reviewVM.get("id"),
        clause: row.get("clauseid"),
        data: reviewVM.get("updateParam")
    }, function(resUpdate){
		leftCol.checkAllConfirmed()
		leftCol.setNextRowToActive()

		self.setConfirmStatus(row)

		unSelect()
    })
}

quizVM.captureHighlighted = function(){
	var self = this

	var highlightClause = viewerVM.get("highlightClause")

    if(highlightClause != ""){

        ((function() {
            var selected = doGetSelection()

            if(viewerVM.highlightValidation(selected) == true){
                var selectedText = captureSelection()
                
                var currentRow = leftCol.get("currentRow")

                currentRow.set("contentFromUser", "Yes")
                currentRow.set("contentFromSystem", selectedText)
                currentRow.set("isIdentified", true)
                currentRow.set("confirmOrTrained", "Confirmed")

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
                    "identified": true,
                    "confirmortrained": currentRow.get("confirmOrTrained")
                }])

                currentRow.set("metadata", reviewVM.get("updateParam"))

                $("#btn-capture").hide()
                currentRow.set("status", true)
                currentRow.set("trainStatus", "Confirmed")

                self.doUpdate(function(){
                    var dizIndex = _.indexOf(currentRow.parent(), _.find(currentRow.parent(), function(v){ return v.uid == currentRow.uid }))

                    if (dizIndex < currentRow.parent().length - 1){
                        currentRow.set("active", false)
                        currentRow.parent()[dizIndex + 1].set("active", true)
                    }
                    unSelect()

                    viewerVM.set("highlightClause", "")
                    reviewVM.refresh(leftCol.getNextRowIndex())
                })
            }
        })())
    }
}

quizVM.doUpdate = function(callback){
	var param = { 
        id: reviewVM.get("id"),
        clause: viewerVM.get("highlightClause"),
        data: reviewVM.get("updateParam")
    }

	ajaxPost("/upload/updatequiz", param, function(resUpdate){
        ajaxPost("gethtmlfilename?id=" + readGetFromUrl().id[0], {}, function(resFile){ 
            var filename = resFile.Data.filenamesource
            reviewVM.set("data", resFile.Data)
            reviewVM.set("subcluster", (resFile.Data.cluster + resFile.Data.subcluster).replace(/\s/g, ''))

            viewerVM.set("highlightClause", "")
			leftCol.checkAllConfirmed();
			viewerVM.highlight()

    		if(callback != undefined) 
    			callback()
        })
    })
}

quizVM.back2Meta = function(e){
	ajaxPost("/upload/updatemetadata", { 
        id: reviewVM.get("id"),
        clause: viewerVM.get("highlightClause"),
        data: reviewVM.get("updateParam"),
        status: ""
    }, function(resUpdate){
        ajaxPost("gethtmlfilename?id=" + readGetFromUrl().id[0], {}, function(resFile){ 
            var filename = resFile.Data.filenamesource
            reviewVM.set("data", resFile.Data)
            reviewVM.set("subcluster", (resFile.Data.cluster + resFile.Data.subcluster).replace(/\s/g, ''))

            e.data.set("status", null)
            leftCol.set("p", "")
            viewerVM.set("highlightClause", "")
    		reviewVM.refresh()
        })
    })
}

quizVM.saveNClose = function(e){
    window.location = "/dashboard/default?p=2";
	// var status = "META"
	// ajaxPost("/upload/updatemetadata", { 
 //        id: reviewVM.get("id"),
 //        clause: reviewVM.get("clause"),
 //        data: reviewVM.get("updateParam"),
 //        status: status
 //    }, function(resUpdate){
 //        ajaxPost("gethtmlfilename?id=" + readGetFromUrl().id[0], {}, function(resFile){ 
 //            var filename = resFile.Data.filenamesource
 //            reviewVM.set("data", resFile.Data)
 //            reviewVM.set("subcluster", (resFile.Data.cluster + resFile.Data.subcluster).replace(/\s/g, ''))

 //            e.data.set("status", null)
 //            leftCol.set("p", status)
 //            reviewVM.set("clause", "")
 //        })
 //    })
 //    leftCol.checkAllConfirmed();
 //    reviewVM.refresh()
}

$(document).ready(function(){
	kendo.bind($("#questionnairePage"), quizVM);
})