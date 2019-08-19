var generalLedger = {}
generalLedger.DateStart = ko.observable(moment().startOf('month').format('YYYY-MM-DD hh:mm'))
generalLedger.DateEnd = ko.observable(new Date)
generalLedger.dataMaster = ko.observableArray([])
generalLedger.dataHead = ko.observableArray([])
generalLedger.TitelFilter = ko.observable(" Show Filter")
generalLedger.DatePageBar = ko.observable()
generalLedger.textSearch = ko.observable("")
generalLedger.dateStart = ko.observable(new Date)
generalLedger.dateEnd = ko.observable(new Date)

generalLedger.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    generalLedger.DatePageBar(page)
}
generalLedger.getData = function (callback) {
    model.Processing(true)
    var s = $('#dateStart').data('kendoDatePicker').value();
    var e = $('#dateEnd').data('kendoDatePicker').value();
    var param = {
        Type: "",
        DateStart: moment(s).format('YYYY-MM-DD'),
        DateEnd: moment(e).format('YYYY-MM-DD'),
        Filter: true,
        TextSearch: generalLedger.textSearch().toUpperCase(),
    }
    ajaxPost('/report/getdatageneralledger', param, function (res) {
        model.Processing(false)
        if (res.IsError) {
            swal("Search Not Found!", res.Message, "warning")
            $('#textSearch').val("")
            return
        }
        generalLedger.dataMaster(res.Data)
        callback()
    }, function () {
        swal("Error!", "Unknown error, please try again", "error")
    })
}
generalLedger.getDataHead = function () {
    var data = [{
        "Acc_Code": 1121,
        "Bank_Account": "USD MANDIRI",
        "Description": "USD Apr 2017",
        "Begining": 0,
        "Transaction": 0,
        "Ending": 0
    }, {
        "Acc_Code": 1122,
        "Bank_Account": "IDR MANDIRI",
        "Description": "IDR Apr 2017",
        "Begining": 0,
        "Transaction": 0,
        "Ending": 0
    }, {
        "Acc_Code": 1123,
        "Bank_Account": "PTC - MANDIRI",
        "Description": "Petty Cash Mandiri Apr 2017",
        "Begining": 0,
        "Transaction": 0,
        "Ending": 0
    }, {
        "Acc_Code": 1110,
        "Bank_Account": "PETTY CASH",
        "Description": "Petty Cash Apr 2017",
        "Begining": 0,
        "Transaction": -80000,
        "Ending": 0
    }]
    generalLedger.dataHead(data)
}

generalLedger.renderheadGrid = function () {
    var data = generalLedger.dataHead();
    if (typeof $('#gridHeadGeneralLedger').data('kendoGrid') !== 'undefined') {
        $('#gridHeadGeneralLedger').data('kendoGrid').setDataSource(new kendo.data.DataSource({
            data: data,
            pageSize: 25,
        }))
        return
    }

    var columns = [{
        field: "Acc_Code",
        title: "Acc_Code"
    }, {
        field: 'Bank_Account',
        title: 'Bank_Account'
    }, {
        field: 'Description',
        title: 'Description'
    }, {
        field: 'Begining',
        title: 'Begining',
        template: function (e) {
            if (e.Begining >= 0) {
                var TotString = kendo.toString(e.Begining, "n");
                return TotString;
            } else {
                var TotminString = kendo.toString(Math.abs(e.Begining), "n");
                return "(" + TotminString + ")";
            }
        }
    }, {
        field: 'Transaction',
        title: 'Transaction',
        template: function (e) {
            if (e.Transaction >= 0) {
                var TotString = kendo.toString(e.Transaction, "n");
                return TotString;
            } else {
                var TotminString = kendo.toString(Math.abs(e.Transaction), "n");
                return "(" + TotminString + ")";
            }
        }
    }, {
        field: 'Ending',
        title: 'Ending',
        template: function (e) {
            var ending = e.Begining + e.Transaction
            if (ending >= 0) {
                var TotString = kendo.toString(ending, "n");
                return TotString;
            } else {
                var TotminString = kendo.toString(Math.abs(ending), "n");
                return "(" + TotminString + ")";
            }
        }
    }]

    $('#gridHeadGeneralLedger').kendoGrid({
        dataSource: data,
        sortable: true,
        scrollable: true,
        columns: columns
    })
}
generalLedger.renderGrid = function () {
    var data = generalLedger.dataMaster();
    var columns = [{
        title: "No.",
        filterable: false,
        template: function (dataItem) {
            var idxs = _.findIndex(data, function (d) {
                return d._id == dataItem._id
            })

            return idxs + 1
        },
        width: 35,
    }, {
        field: 'DateStr',
        title: 'Date',
        width: 80,
        template: function(e){
                var date = moment(e.DateStr).format("DD-MMM-YYYY")
                if(date=="01-Jan-0001"){
                    date = ""
                }
                return date
            }
    }, {
        field: 'DocumentNumber',
        title: 'Document Number',
        width: 100
    }, {
        field: 'Acc_Code',
        title: 'Acc Code',
        width: 80
    }, {
        field: 'Acc_Name',
        title: 'Account Name',
        width: 100
    }, {
        field: 'Department',
        title: 'Department',
        width: 70,
    },  {
        field: 'SalesName',
        title: 'Sales',
        width: 70,
    }, {
        field: 'Description',
        title: 'Description',
        width: 130
    }, {
        field: 'IdJournal',
        title: 'Reff',
        width: 100,
        template: function (d) {
            return '<a onclick="generalLedger.gridModal(\'' + d.IdJournal + '\')" data-target="#gridModal" data-toggle="modal">' + d.IdJournal + '</a>'
        },
        attributes: {
            "class": "colorText",
        },
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>Total:</div>"
    }, {
        field: 'Debet',
        title: 'Debit',
        width: 100,
        template: "#=ChangeToRupiah(Debet)#",
        aggregates: ["sum"],
        attributes: {
            style: "text-align:right;"
        },
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>"
    }, {
        field: 'Credit',
        title: 'Credit',
        width: 100,
        template: "#=ChangeToRupiah(Credit)#",
        aggregates: ["sum"],
        attributes: {
            style: "text-align:right;"
        },
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>"
    }]
    $('#gridGeneralLedger').html("")
    $('#gridGeneralLedger').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
                field: "Debet",
                aggregate: "sum"
            }, {
                field: "Credit",
                aggregate: "sum"
            }],
        },
        excel: {
            fileName: "report-generalledger.xlsx"
        },
        excelExport: function (dataItem) {
            var rows = dataItem.workbook.sheets[0].rows;
            for (var ri = 0; ri < rows.length; ri++) {
                var row = rows[ri];
                if (row.type == "group-footer" || row.type == "footer") {
                    for (var ci = 0; ci < row.cells.length; ci++) {
                        var cell = row.cells[ci];
                        if (cell.value) {
                            // Use jQuery.fn.text to remove the HTML and get only the text
                            var text = $(cell.value).text();
                            if (text != "Total:") {
                                var expression = /^\([\d,\.]*\)$/;
                                if (text.match(expression)) {
                                    //It matched - strip out parentheses and append - at front
                                    var val = '-' + text.replace(/[\(\)]/g, '')
                                } else {
                                    var val = text
                                }
                                cell.value = parseFloat(val.split(",").join(""))
                                cell.format = "#,##0.00_);(#,##0.00);0.00;"
                                // Set the alignment
                                cell.hAlign = "right";
                            } else {
                                cell.value = text
                            }
                        }
                    }
                }
                if (row.type == "data") {
                    if (ri > 0) {
                        for (var ci = 0; ci < row.cells.length; ci++) {
                            var cell = row.cells[ci];
                            if (ci > 7) {
                                cell.value = parseFloat(cell.value)
                                cell.format = "#,##0.00_);(#,##0.00);0.00;"
                                // Set the alignment
                                cell.hAlign = "right";
                            }
                        }
                    }
                }
            }
        },
        height: 500,
        sortable: true,
        scrollable: true,
        columns: columns,

    })
}

generalLedger.gridModal = function (ids) {
    var data = _.filter(generalLedger.dataMaster(), {
        IdJournal: ids
    });
    var grid = new kendo.data.DataSource({
        data: data,
        pageSize: 25,
    })
    var columns = [{
        title: "No.",
        filterable: false,
        template: function (dataItem) {
            var idx = _.findIndex(data, function (d) {
                return d._id == dataItem._id
            })

            return idx + 1
        },
        width: 30
    }, {
        field: 'DateStr',
        title: 'Date',
        width: 100,
         template: function(e){
                var date = moment(e.DateStr).format("DD-MMM-YYYY")
                if(date=="01-Jan-0001"){
                    date = ""
                }
                return date
            }
    }, {
        field: 'DocumentNumber',
        title: 'DocumentNumber',
        width: 100
    }, {
        field: 'Acc_Code',
        title: 'Acc Code',
        width: 100
    },

    {
        field: 'Acc_Name',
        title: 'Acc Name',
        width: 100
    }, {
        field: 'Description',
        title: 'Description',
        width: 100
    }, {
        field: 'Debet',
        title: 'Debit',
        width: 100,
        template: "#=ChangeToRupiah(Debet)#",
        aggregates: ["sum"],
        attributes: {
            style: "text-align:right;"
        },
        footerTemplate: "#= sum #",
    }, {
        field: 'Credit',
        title: 'Credit',
        width: 100,
        aggregates: ["sum"],
        attributes: {
            style: "text-align:right;"
        },
        template: "#=ChangeToRupiah(Credit)#",
        footerTemplate: "#= sum #",
    }

    ]

    $('#gridReffGeneralLedger').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
                field: "Debet",
                aggregate: "sum"
            }, {
                field: "Credit",
                aggregate: "sum"
            }, ],
        },
        scrollable: true,
        columns: columns
    })
}

generalLedger.refreshGrid = function () {
    generalLedger.getData(function () {
        generalLedger.renderGrid()
    })
}

generalLedger.search = function() {
    // closing.refresh(true)
    generalLedger.getData(function () {
        generalLedger.renderGrid()
    })
 }

generalLedger.setDate = function () {
    var datepicker = $("#dateStart").data("kendoDatePicker");
    var now = new Date();
    var years = moment(now).format("YYYY")
    var Month = moment(now).format("MM")
    var newDate = years + "-" + Month + "-" + "01"
    datepicker.value(new Date(newDate))
}
generalLedger.ExportToPdf = function(){
    var Data = $("#gridGeneralLedger").data("kendoGrid").dataSource.options.data
    Data = _.sortBy(Data, function(o){return o.PostingDate})
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var url = "/report/exportpdfgeneralledger"
    var param = {
        // Filter : mCoa.DateFilter(),
        DateStart: dateStart,
        DateEnd: dateEnd,
        Data : Data,
    }
    model.Processing(true)
    ajaxPost(url, param, function(e){
        window.open('/res/docs/report/pdf/' + e.Data, '_blank');
        model.Processing(false)
    })
}
generalLedger.dowloadExcel = function () {
    $("#gridGeneralLedger").getKendoGrid().saveAsExcel();
}

generalLedger.onChangeDateStart = function(val){
    if (val.getTime()>generalLedger.DateEnd().getTime()){
        generalLedger.DateEnd(val)
    }
}

generalLedger.init = function () {
    // generalLedger.setDate()
    generalLedger.getData(function () {
        generalLedger.renderGrid()
    })
    generalLedger.getDataHead()
    generalLedger.renderheadGrid()
    generalLedger.getDateNow()
}

$(document).ready(function () {
    generalLedger.init();
})