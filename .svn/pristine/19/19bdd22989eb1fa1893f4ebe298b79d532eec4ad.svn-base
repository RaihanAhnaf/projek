var reviewVM = kendo.observable({
	data: null,
	subcluster: "",
	id: "",
    confidence: "",
    updateParam: [],
    statussummary: "",
    refresh: function(destinationRowIndex){
    	var self = this

		var params = readGetFromUrl();
		ajaxPost("/countrymaster/getcountry", {}, function(resCountry){
			var cm = resCountry.Res.Data

			ajaxPost("/dashboard/gethtmlfilename?id=" + params.id[0], {}, function(resFile){
				var dizFile = resFile.Data

				var filename = dizFile.filenamesource
				var subcluster = (dizFile.subcluster).replace(/\s/g, '')

				var metaActive = 0

				var getGoverningLawName = function(){
					if(subcluster == "externalmsa" || subcluster == "externalsla"){
						return "Governing Law"
					} else {
						return "Receiving Entities"
					}
				}

				var arrMeta = _.map(dizFile.metadata, function(meta, key){
					var ret = {
						clauseid: key,
						clauselabel: (key == "governinglaw") ? getGoverningLawName() : metadataMapping[subcluster][key].name,
						type: metadataMapping[subcluster][key].type, //type
						contentFromSystem: (meta.length > 0) ? meta[meta.length - 1].contentfromsystem : "",
						contentFromUser: (function(data){
							if(metadataMapping[subcluster][key].type == 'date')
								return moment(data).toDate()
							else
								return data
						})((meta.length > 0) ? meta[meta.length - 1].contentfromuser : ""),
						metadata: meta,
						isIdentified: (meta.length > 0) ? meta[meta.length - 1].identified : false, 
						confirmOrTrained: (meta.length > 0) ? meta[meta.length - 1].confirmortrained : "",
						isConfirmed: false,
						active: false
					}

					ret.list = new kendo.data.DataSource({
				    	data: cm
				    })

					return ret
				});
				metadataVM.set("arraj", arrMeta)

				var quizActive = 0
				if(subcluster == "externalmsa" || subcluster == "externalsla"){
					var arrQuiz = _.map(dizFile.questionnaire, function(quiz, key){
						var ret = {
							clauseid: key,
							clauselabel: questionnaireMapping[subcluster][key].name,
							sentence: questionnaireMapping[subcluster][key].sentence, //type
							contentFromSystem: (quiz.length > 0) ? quiz[quiz.length - 1].contentfromsystem : "",
							contentFromUser: (quiz.length > 0) ? quiz[quiz.length - 1].contentfromuser : "",
							metadata: quiz,
							isIdentified: (quiz.length > 0) ? quiz[quiz.length - 1].identified : false, 
							confirmOrTrained: (quiz.length > 0) ? quiz[quiz.length - 1].confirmortrained : "",
							active: false
						}

						return ret
					});
					quizVM.set("arraj", arrQuiz)
				}

				self.set("data", dizFile)
				self.set("id", params.id[0])
				self.set("subcluster", subcluster)
				self.set("statussummary", dizFile.statussummary)

				//subcluster type
				viewerVM.set("contractClassification", toTitleCase(dizFile.cluster))
				viewerVM.set("contractType", dizFile.subcluster.toUpperCase().replace(dizFile.cluster.toUpperCase(), ""))
				viewerVM.set("filename", filename)
				viewerVM.set("uploadTime", moment(dizFile.uploadtime).format("DD MMM YYYY"))
				viewerVM.set("pdfurl", "/files/UnprocessDoc/" + filename)
				viewerVM.set("htmlurl", "/files/Html/" + filename)

				leftCol.set("quizBarVisible", (self.get("subcluster") == "externalmsa" || self.get("subcluster") == "externalsla"))

				// var urlviewer = "{{BaseUrl}}res/pdfjs/web/viewer.html?file=/files/Pdf/" + filename + ".pdf"
				var urlviewer = "/files/Html/" + filename

				viewerVM.set("dizsrc", urlviewer)
				document.getElementById('pdfframe').src = viewerVM.get("dizsrc")
				
				leftCol.set("p", dizFile.status)
				if(leftCol.get("p") == ""){
					leftCol.set("currentVM", metadataVM)
				} else if(leftCol.get("p") == "META"){
					leftCol.set("currentVM", quizVM)
				}
				leftCol.get("currentVM").init()

				if( ! destinationRowIndex)
					leftCol.setRowToActive(leftCol.get("currentVM").get("arraj")[0])
				else
					if(leftCol.get("currentVM").get("arraj")[destinationRowIndex] != undefined)
						leftCol.setRowToActive(leftCol.get("currentVM").get("arraj")[destinationRowIndex])
					else
						leftCol.setRowToActive(leftCol.get("currentVM").get("arraj")[0])

				resizeWidth()
			})
		})
	}
})

//-------------------------------------------------------------------------------------------------------------------------

var viewerVM = kendo.observable({
	dizsrc: "",
	contractClassification: "",
	contractType: "",
	filename: "",
	uploadTime: "",
	pdfurl: "",
	htmlurl: "",
	findString: function(str) {
	  if (parseInt(navigator.appVersion) < 4) return;
	  
	  var strFound;
	  
	  if (document.getElementById("pdfframe").contentWindow.find) {

	   	var self = document.getElementById("pdfframe").contentWindow

	    strFound = self.find(str);
	    if ( ! strFound) {
	      strFound = self.find(str, 0, 1);
	      while (self.find(str, 0, 1)) continue;
	    }
	  } else if (navigator.appName == "Opera") {
	    alert ("Opera browsers not supported, sorry...")
	    return;
	  }
	},
	keypressed: function(e){
	    if (e.which == 13) {
	      	viewerVM.findString(e.target.value)
	    } 
	},
	highlightValidation: function(selected){
		if(parseInt(selected.startElementId) > parseInt(selected.endElementId) || selected.endElementId == null){
			swal("Error!", "Selection invalid, please reselect.", "error");
			
			$("#modal-review").modal("hide")
			unSelect()

			viewerVM.set("highlightClause", leftCol.get("currentRow").clauseid)

			leftCol.checkAllConfirmed();
			return false
		} else {
			return true
		}
	},
	highlight: function() {
		document.getElementById('pdfframe').src = document.getElementById('pdfframe').src

		$('#pdfframe').load(function(){
			$(document.getElementById("pdfframe").contentWindow.document.body.querySelector('span.highlight')).replaceWith(function(){
		        return $(this).text()
		    });

			var currentRow = leftCol.get("currentRow")
			var foundElemsOffset = []

			_.each(currentRow.get("metadata"), function(meta){
				var start = meta["dataeaciitstart"] != undefined ? meta["dataeaciitstart"] : meta["data-eaciit-start"]
				var end = meta["dataeaciitend"] != undefined ? meta["dataeaciitend"] : meta["data-eaciit-end"]

		        var grabbedElements = [];
		        _.each(document.getElementById("pdfframe").contentWindow.document.body.getElementsByTagName("*"), function(v) {
		            if( ! $(v).is(":visible")) return;
		            
		            if(parseInt(v.getAttribute("data-eaciit-id")) >= start 
		                && parseInt(v.getAttribute("data-eaciit-id")) <= end 
		                && v.hasAttribute("data-eaciit-start") 
		                && v.hasAttribute("data-eaciit-end")
		                ){
		                grabbedElements.push(v);
		            }
		        });

		        _.each(grabbedElements, function(grabbed){
					(function(id, text) {
					    var foundElem = (function(attribute, value) {
							var all = document.getElementById("pdfframe").contentWindow.document.body.getElementsByTagName('*');
							for (var i = 0; i < all.length; i++){
								if (all[i].getAttribute(attribute) == value) { 
									return all[i]; 
								}
							}
						})("data-eaciit-id", id)

						if(foundElem != undefined){
							var isinya = foundElem.innerHTML
							isinya = $('<textarea />').html(isinya).text(); //removes html entities
							isinya = isinya.replace(/\n/g, " ");			//removes newlines
							isinya = isinya.replace(/\s\s+/g, ' ');			//removes double spaces

						    var indexWithElems = isinya.indexOf(text);
						    if (indexWithElems >= 0){ 
						    	var div = document.createElement("div");
								div.innerHTML = isinya;

						        isinya = isinya.substring(0, indexWithElems) + "<span class='highlight' style='background-color: yellow'>" + isinya.substring(indexWithElems, indexWithElems + text.length) + "</span>" + isinya.substring(indexWithElems + text.length);
						        foundElem.innerHTML = isinya 
						    } else {
						    	foundElem.style.backgroundColor = 'yellow';
						    }

							addClass(foundElem, "terhighlight")

						    elemOffset = (function(element) {
							    var top = 0, left = 0;
							    do {
							        top += element.offsetTop  || 0;
							        left += element.offsetLeft || 0;
							        element = element.offsetParent;
							    } while(element);

							    return {
							        top: top,
							        left: left
							    };
							})(foundElem)

						    foundElemsOffset.push(elemOffset.top)
						}
					})(grabbed.getAttribute("data-eaciit-id"), meta.contentfromsystem)
		        })
			});
			document.getElementById("pdfframe").contentWindow.scrollTo(0, _.min(foundElemsOffset) - 200)
		})
	}
})

//-------------------------------------------------------------------------------------------------------------------------

var leftCol = kendo.observable({
	p: "",
	currentVM: null,
	quizBarVisible: false,
	currentRow: null,
	areWeTakingTheQuiz: function(){
		if(reviewVM.get("subcluster") == "externalmsa" || reviewVM.get("subcluster") == "externalsla"){
			return true
		} else {
			return false
		}
	},
	isMetadataPage: function(){ 
		return (this.get("p") == "" || ( ! this.areWeTakingTheQuiz())) 
	},
	isQuizPage: function(){ 
		return this.get("p") == "META" 
	},
	nextStep: function(){
		if (this.get("p") == "")
			this.set("p", "META")
	},
	setRowToActive: function(row){
		if(row.get("active") == false){
			var currentVM = this.get("currentVM")

			var activeRow = _.find(currentVM.get("arraj"), function(meta){
				return meta.get("active") == true
			});

			if(activeRow != undefined)
				activeRow.set("active", false)

			row.set("active", true)
			this.set("currentRow", row)

			modalTraining.set("currentClause", row.clauselabel)
			modalReview.set("currentClause", row.clauselabel)

			viewerVM.set("highlightClause", "")
			viewerVM.highlight()
		}
	},
	getNextRowIndex: function(){
		var currentVMArray = this.get("currentVM").get("arraj")
		var currentRow = this.get("currentRow")

		return (_.indexOf(currentVMArray, _.find(currentVMArray, function(v){ 
			return v.uid == currentRow.uid 
		})) + 1)
	},
	getNextRow: function(){
		var currentVMArray = this.get("currentVM").get("arraj")
		return currentVMArray[this.getNextRowIndex()]
	},
	setNextRowToActive: function(){
		var nextRow = this.getNextRow()
		
		if(nextRow != undefined)
			this.setRowToActive(nextRow)
		else
			viewerVM.highlight()
	},
	checkAllConfirmed: function(){
		var currentVM = this.get("currentVM")

		var countIdentified = _.countBy(currentVM.get("arraj"), function(v) {
		  return v.get("isIdentified") == true ? "yes" : "no"
		});

		if(countIdentified.no == undefined)
			currentVM.set("allConfirmed", true)

		var countConfirmOrTrained = _.countBy(currentVM.get("arraj"), function(v) {
		  return v.get("confirmOrTrained") != "" ? "yes" : "no"
		});

		if(countConfirmOrTrained.no == undefined)
			currentVM.set("allConfirmOrTrained", true)
	},
	trainRow: function(row){
		$("#modal-training").modal("show")

		viewerVM.set("highlightClause", row.clauseid)
	},
	setRowNotFound: function(row){
		var self = this

		var currentVM = self.get("currentVM")

		row.set("isIdentified", true)
		row.set("confirmOrTrained", "Not Found")

	    reviewVM.set("updateParam", [{
	        "confidence": 1.0,
	        "clause": row.get("clauseid"),
	        "contentfromuser": "",
	        "contentfromsystem": "",
	        "common-ancestor-id": "",
	        "data-eaciit-id": "",
	        "start-offset": "",
	        "end-offset": "",
	        "data-eaciit-start": "",
	        "data-eaciit-end": "",
	        "start-offset": "",
	        "identified": row.get("isIdentified"),
	        "confirmortrained": row.get("confirmOrTrained")
	    }])

	    ajaxPost(self.get("currentVM").url.update, { 
	        id: reviewVM.get("id"),
	        clause: row.get("clauseid"),
	        data: reviewVM.get("updateParam")
	    }, function(resUpdate){
			self.checkAllConfirmed()
			// self.setNextRowToActive()

			currentVM.setConfirmStatus(row)

			unSelect()
        	
        	reviewVM.refresh(leftCol.getNextRowIndex())
	    })
	}
})

//-------------------------------------------------------------------------------------------------------------------------

var modalTraining = kendo.observable({
	currentClause: "",
	cancel: function(){
		$("#modal-training").modal("hide")
		viewerVM.set("highlightClause", "")
	},
	doSelect: function(){
		$("#modal-training").modal("hide")
		unSelect()
	}
})

//-------------------------------------------------------------------------------------------------------------------------

var modalReview = kendo.observable({
	currentClause: "",
	select: function(){
		$("#modal-review").modal("hide")
		unSelect()
	},
	confirm: function(){
		$("#modal-training").modal("hide")

		leftCol.get("currentVM").captureHighlighted()
		$("#modal-review").modal("hide")
	}
})

//-------------------------------------------------------------------------------------------------------------------------

var findByAttributeValue = function(attribute, value) {
	var all = document.getElementById("pdfframe").contentWindow.document.body.getElementsByTagName('*');
	for (var i = 0; i < all.length; i++){
		if (all[i].getAttribute(attribute) == value) { 
			return all[i]; 
		}
	}
}

var countryMaster = kendo.observable({
	isSet: function(){
		return countryMaster.get("arr").length > 0
	},
	arr: []
})

var resizeWidth = function(){
    $("#appS > ul > li").width(($("#appS > ul").width() / 2) - 60)
    $('<style>li.panah-biasa:after { content: ""; position: absolute; left: '+(($("#appS > ul").width() / 2) - 30)+'px; width: 0; height: 0; border-top: 20px solid transparent; border-left: 20px solid #6a80a2; border-bottom: 20px solid transparent; margin: -10px 90px 0 10px; }#appS li.active:after { border-left: 20px solid #313d50 !important; }</style>').appendTo("#appS > ul > li")

    $(".col-sisa-inputan").width(191)
    $(".col-inputan").width($("#col-parent-inputan").width() - 191)
    $(".col-quiz").width($("#col-parent-quiz").width() - 191)
}

$(window).resize(resizeWidth)

$(document).ready(function () {
	kendo.bind($("#review"), reviewVM);
	kendo.bind($("#rightCol"), viewerVM)
	kendo.bind($("#leftCol"), leftCol);
	kendo.bind($("#modal-training"), modalTraining);
	kendo.bind($("#modal-review"), modalReview);

	ajaxPost("/countrymaster/getcountry", {}, function(res){
		if ( ! countryMaster.isSet()){
			countryMaster.set("arr", res.Res.Data)
		}
	})

	setTimeout(resizeWidth, 500)

	reviewVM.refresh()
});