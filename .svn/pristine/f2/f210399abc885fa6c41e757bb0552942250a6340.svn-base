var chartDates = kendo.observable({
	classification: new Date(),
	type: new Date(),
    chartmeta1: false,
    chartmeta2: false,
    chartmeta3: false,
    chartquiz1: false,
    chartquiz2: false,
    chartquiz3: false,
    metadata1: function(e){
        addchartModal.set("chartnumber", 1)
        addchartModal.init("metadata")
    },
    metadata2: function(){
        addchartModal.set("chartnumber", 2)
        addchartModal.init("metadata")
    },
    metadata3: function(){
        addchartModal.set("chartnumber", 3)
        addchartModal.init("metadata")
    },
    questionnaire1: function(){
        addchartModal.set("chartnumber", 1)
        addchartModal.init("questionnaire")
    },
    questionnaire2: function(){
        addchartModal.set("chartnumber", 2)
        addchartModal.init("questionnaire")
    },
    questionnaire3: function(){
        addchartModal.set("chartnumber", 3)
        addchartModal.init("questionnaire")
    }
})

var addchartModal = kendo.observable({
    chartrace: "",
    chartnumber: null,
    dateVal: new Date(),
    classificationVal: null,
    typeVal: null,
    metadataVal: null,
    classifications: [],
    types: [],
    metadatas: [],
    reset: function(){
        var self = this

        $('#dp-modal').daterangepicker({
            singleDatePicker: true,
            maxDate: new Date()
        }, function(start, end, label) {
            addchartModal.set("dateVal", start.toDate())
        });

        addchartModal.set("dateVal", new Date())
        addchartModal.set("classificationVal", null)
        addchartModal.set("typeVal", null)
        addchartModal.set("metadataVal", null)

        addchartModal.set("classifications", [])
        addchartModal.set("types", [])
        addchartModal.set("metadatas", [])
    },
    setClassifications: function(){
        var self = this

        ajaxPost("/refmaster/getdatamaster", {
            category: "DOC-CLUSTER"
        }, function(res){
            self.set("classifications", res.Res.Data)
        })
    },
    setTypes: function(){
        var self = this

        ajaxPost("/refmaster/getdatamaster", {
            category: "DOC-SUBCLUSTER"
        }, function(res){
            var types = _.filter(res.Res.Data, function(value, key, list){
                if(self.get("classificationVal") != null){
                    return value.Code.indexOf(self.get("classificationVal").Code) != -1;
                } else
                    return true
            });

            self.set("types", types)
        })
    },
    init: function(race){
        var self = this

        self.reset()

        self.setClassifications()
        self.setTypes()

        switch(race){
            case "metadata": self.set("metadatas", (function(){
                self.set("chartrace", "meta")
                
                return _.map(metadataRows, function(v, k){
                    return {
                        "Code": k,
                        "Description": v.name
                    }
                })
            })()); break;

            case "questionnaire": self.set("metadatas", (function(){
                self.set("chartrace", "quiz")
                
                return _.map(questionnaireRows, function(v, k){
                    return {
                        "Code": k,
                        "Description": v.name
                    }
                })
            })()); break;
        }
    },
    add: function(){
        var self = this;

        var chartfound = "chart" + self.get("chartrace") + self.get("chartnumber")
        chartDates.set(chartfound, true)

        $("#" + chartfound).kendoChart({
            dataSource: {
                serverPaging: true,
                pageSize: 100,
                transport: {
                    read: function(o){
                        ajaxPost("/analytics/generatedatachartbymetadata", {
                            asofdate: self.get("dateVal"),
                            cluster: self.get("classificationVal").Code,
                            subcluster: self.get("typeVal").Code,
                            clause: self.get("metadataVal").Code,
                        }, function(res){
                            var data = res.Data

                            _.each(data, function(value){
                                value.External = value.External != null ? parseFloat(value.External).toFixed() : null;
                                value.Internal = value.Internal != null ? parseFloat(value.Internal).toFixed() : null;
                            });

                            data.push({
                                Category: "",
                                Confidence: null
                            })

                            o.success(data);
                        })
                    }
                }
            },
            title: {
                text: toTitleCase((function(){ 
                    switch(self.get("chartrace")){
                        case "meta": return "Metadata"; break;
                        case "quiz": return "Questionnaire"; break;
                    }
                })())
            },
            legend: {
               position: "bottom"
            },
            seriesDefaults: {
                type: "line",
                style: "smooth"
            },
            series: [{
                field: "Confidence",
                name: "Confidence",
                color: '#0099ff'
            }],
            valueAxis: {
                title: {
                    text: "Classification Accuracy",
                    font:"12px Helvetica Neue, Helvetica, Arial, sans-serif",
                },
                labels: {
                    format: "{0}%"
                },
                line: {
                    visible: true
                },
                majorGridLines: {
                    visible: false
                },
                axisCrossingValue: -10,
                max: 105
            },
            categoryAxis: {
                field: "Category",
                title: {
                    text: "No. of Labeled Docs",
                    font:"12px Helvetica Neue, Helvetica, Arial, sans-serif",
                },
                majorGridLines: {
                    visible: false
                },
                labels: {
                    rotation: "auto"
                },
                justified: true
            },
            tooltip: { 
                visible: true, 
                format: "{0}%" 
            }
        })
    }
})

var createChart = function() {
   $("#chart-classification").kendoChart({
        dataSource: {
            serverPaging: true,
            pageSize: 100,
            transport: {
                read: function(o){
                    ajaxPost("/analytics/generatedatachartbycluster", {
                        asofdate: chartDates.get("classification"),
                    }, function(res){
                    	var data = res.Data

                        _.each(data, function(value){
                            value.External = value.External != null ? parseFloat(value.External).toFixed() : null;
                            value.Internal = value.Internal != null ? parseFloat(value.Internal).toFixed() : null;
                        });

                        data.push({
                        	Category: "",
							External: null,
							Internal: null
                        })

                        o.success(data);
                    })
                }
            }
        },
       	title: {
           	text: "Classification Accuracy: Internal v External"
       	},
       	legend: {
           position: "bottom"
       	},
       	seriesDefaults: {
			type: "line",
			style: "smooth"
       	},
       	series: [{
       		field: "Internal",
			name: "Internal",
			color: '#0099ff'
       	},
       	{
       		field: "External",
			name: "External",
			color: '#2db245'
       	}],
       	valueAxis: {
          	title: {
	            text: "Classification Accuracy",
	            font:"12px Helvetica Neue, Helvetica, Arial, sans-serif",
	        },
            labels: {
                format: "{0}%"
            },
           	line: {
               	visible: true
           	},
         	majorGridLines: {
               	visible: false
           	},
           	axisCrossingValue: -10,
           	max: 105
       	},
       	categoryAxis: {
       		field: "Category",
         	title: {
				text: "No. of Labeled Docs",
				font:"12px Helvetica Neue, Helvetica, Arial, sans-serif",
           	},
           	majorGridLines: {
               	visible: false
           	},
           	labels: {
               	rotation: "auto"
           	},
           	justified: true
		},
		tooltip: { 
			visible: true, 
			format: "{0}%" 
		}
	});

	$("#chart-type").kendoChart({
        dataSource: {
            serverPaging: true,
            pageSize: 100,
            transport: {
                read: function(o){
                    ajaxPost("/analytics/generatedatachartbysubcluster", {
                        asofdate: chartDates.get("type"),
                    }, function(res){
                        var data = res.Data

                        _.each(data, function(value){
                            value.Msa = value.Msa != null ? parseFloat(value.Msa).toFixed() : null;
                            value.Sla = value.Sla != null ? parseFloat(value.Sla).toFixed() : null;
                        });

                        data.push({
                            Category: "",
                            Msa: null,
                            Sla: null
                        })
                        
                        o.success(data);
                    })
                }
            }
        },
       	title: {
           	text: "Model Accuracy (Classifying MSA v SLA v RRP SOC v Addendum"
       	},
       	legend: {
           	position: "bottom"
       	},
       	seriesDefaults: {
           	type: "line",
           	style: "smooth"
       	},
       	series: [{
       		field: "Msa",
			name: "MSA",
			color: '#0099ff'
       	},
       	{
       		field: "Sla",
			name: "SLA",
			color: '#2db245'
       	}],
       	valueAxis: {
          	title: {
	            text: "Classification Accuracy",
	            font:"12px Helvetica Neue, Helvetica, Arial, sans-serif",
	        },
            labels: {
                format: "{0}%"
            },
           	line: {
               	visible: true
           	},
         	majorGridLines: {
               	visible: false
           	},
           	axisCrossingValue: -10,
           	max: 105
       	},
       	categoryAxis: {
       		field: "Category",
         	title: {
				text: "No. of Labeled Docs",
				font:"12px Helvetica Neue, Helvetica, Arial, sans-serif",
           	},
           	majorGridLines: {
               	visible: false
           	},
           	labels: {
               	rotation: "auto"
           	},
           	justified: true
		},
		tooltip: { 
			visible: true, 
			format: "{0}%" 
		}
   	});
}

$('#dp-classification').kendoDatePicker({
    disableDates: function (date) {
        return date > new Date();
    },
    change: function(){
    	chartDates.set("classification", this.value())
    	$("#chart-classification").data("kendoChart").dataSource.read(); 
    	$("#chart-classification").data("kendoChart").refresh();
    }
});
$('#dp-type').kendoDatePicker({
    value: new Date(),
    disableDates: function (date) {
        return date > new Date();
    },
    change: function(){
    	chartDates.set("type", this.value())
    	$("#chart-type").data("kendoChart").dataSource.read(); 
    	$("#chart-type").data("kendoChart").refresh();
    }
});

$(document).ready(function(){
	kendo.bind($("#filtered.bydate"), chartDates)
    kendo.bind($("#addchart-modal"), addchartModal)

    $("#classifications").kendoDropDownList({
        dataValueField: 'Code', 
        dataTextField: 'Name',
        optionLabel: "Select classification..."
    })

    $("#types").kendoDropDownList({
        dataValueField: 'Code', 
        dataTextField: 'Description',
        optionLabel: "Select classification..."
    })

    $("#metadatas").kendoDropDownList({
        dataValueField: 'Code', 
        dataTextField: 'Description',
        optionLabel: "Select classification..."
    })

	createChart()
});
$(document).bind(".chart", createChart);