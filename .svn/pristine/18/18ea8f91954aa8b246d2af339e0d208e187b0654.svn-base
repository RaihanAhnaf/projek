model.Processing(false)
var ledger = {
    TitelFilter: ko.observable(" Hide Filter")
}
ledger.DateStart = ko.observable(moment().startOf('month').format('YYYY-MM-DD hh:mm'))
ledger.DateEnd = ko.observable(new Date)
ledger.koFilterIsActive = ko.observable(false)
ledger.koDataAccount = ko.observableArray([])
ledger.koDataLedger = ko.observableArray([])
ledger.koAccountName = ko.observable("")
ledger.koAccountCode = ko.observable("")
ledger.koParentCode = ko.observable()
ledger.koBeginingValue = ko.observable()
ledger.DatePageBar = ko.observable()
ledger.textSearch = ko.observable("")

ledger.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    ledger.DatePageBar(page)
}

ledger.getAccount = function () {
    var url = "/master/getdatacoa"
    var param = {}
    model.Processing(true)
    ajaxPost(url, param, function (res) {
        model.Processing(false)
        for (i in res.Data) {
            res.Data[i].ACC_Code = res.Data[i].ACC_Code + ""
            res.Data[i].CodeName = res.Data[i].ACC_Code + "-" + res.Data[i].Account_Name
        }
        ledger.koDataAccount(res.Data)
        ledger.renderDropdownAccount(res.Data)
    })
}
ledger.renderDropdownAccount = function (data) {
    $("#DropDownAccountCode").kendoDropDownList({
        dataTextField: "CodeName",
        dataValueField: "ACC_Code",
        dataSource: ledger.koDataAccount(),
        filter: 'contains',
        index: 0,
        optionLabel: "---Select Account---",
        change: function (e) {
            var model = this.dataItem();
            ledger.koAccountName(model.Account_Name)
            ledger.koAccountCode(model.ACC_Code)
            ledger.koParentCode(model.Main_Acc_Code)
        }
    });
}
ledger.getDataLedger = function (callback) {
    model.Processing(true)
    var url = "/report/getdataledger"
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var param = {
        DateStart: dateStart,
        DateEnd: dateEnd,
        Accountcode: ledger.koAccountCode(),
        ParentCode: ledger.koParentCode(),
        Filter: true,
        TextSearch: ledger.textSearch().toUpperCase(),
    }
    ajaxPost(url, param, function (res) {
        var Begining = res.Begining
        var Data = res.Data
        Data = _.sortBy(Data, [function (o) {
            return o.PostingDate
        }])
        var SaldoCalculate = 0.0
        for (i in Data) {
            if (i == 0) {
                var begin = Begining + Data[i].Debet
                Data[i].Saldo = begin - Data[i].Credit
                SaldoCalculate = Data[i].Saldo
            }
            if (i > 0) {
                var begin = SaldoCalculate + Data[i].Debet
                Data[i].Saldo = begin - Data[i].Credit
                SaldoCalculate = Data[i].Saldo
            }
        }
        ledger.koDataLedger(Data)
        ledger.koBeginingValue(ChangeToRupiah(Begining))
        model.Processing(false)
        callback()
    })
}
ledger.renderGrid = function () {
    var data = ledger.koDataLedger()
    var columns = [{
        field: 'PostingDate',
        title: 'Posting Date',
        template: function (e) {
            return moment.utc(e.PostingDate).format("DD-MMM-YYYY")
        },
        width: 100,
    }, {
        field: 'Acc_Code',
        title: 'Acc Code',
        width: 80
    }, {
        field: 'Acc_Name',
        title: 'Account Name',
        width: 100
    }, {
        field: 'DocumentNumber',
        title: 'Document Number',
        width: 80,
    }, {
        field: 'Department',
        title: 'Department',
        width: 70,
    },{
        field: 'SalesName',
        title: 'Sales',
        width: 70,
    },{
        field: 'Description',
        title: 'Description',
        width: 90,
        footerTemplate: "<div style='text-align:center; font-size: 15px;'>Total:</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Debet',
        title: 'Debit',
        width: 90,
        template: "#=ChangeToRupiah(Debet)#",
        footerTemplate: "<div style='text-align:center; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Credit',
        title: 'Credit',
        width: 90,
        template: "#=ChangeToRupiah(Credit)#",
        footerTemplate: "<div style='text-align:center; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Saldo',
        title: 'Saldo',
        width: 90,
        template: "#=ChangeToRupiah(Saldo)#",
        footerTemplate: "<div style='text-align:center; font-size: 15px;'>#= ledger.saldoCalculation() #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, ]
    $('#gridLedger').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
                field: "Debet",
                aggregate: "sum"
            }, {
                field: "Credit",
                aggregate: "sum"
            }, {
                field: "Saldo",
                aggregate: "sum"
            }],
        },
        excel: {
            fileName: "report-ledger.xlsx"
        },
        excelExport: function (e) {
            var rows = e.workbook.sheets[0].rows;
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
                                cell.value = parseFloat(val.split(",").join("")) // this different
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
                            if (ci == 0) {
                                cell.value = moment(cell.value).format("DD MMM YYYY")
                            }
                            if (ci > 6) {
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
        columns: columns
    })
}

ledger.saldoCalculation = function () {
    var data = $('#gridLedger').data('kendoGrid').dataSource.options.data
    var lastData = _.last(data)
    if (lastData != undefined) {
        var Saldo = lastData.Saldo
    } else {
        var Saldo = 0
    }

    return ChangeToRupiah(Saldo)
}
ledger.refreshGrid = function () {
    if (ledger.koAccountCode() == "" || ledger.koAccountName() == "") {
        return swal("Warning!", "Please Select COA Account!", "warning")
    }
    ledger.getDataLedger(function () {
        ledger.renderGrid()
    })
}
ledger.ExportToPdf = function(){
    var Data = $("#gridLedger").data("kendoGrid").dataSource.options.data
    Data = _.sortBy(Data, function(o){return o.PostingDate})
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var url = "/report/exportpdfledger"
    var param = {
        Begining : ledger.koBeginingValue(),
        AccName : ledger.koAccountName(),
        AccCode : ledger.koAccountCode(),
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
ledger.search = function() {
    ledger.refreshGrid(true)
   ledger.getDataLedger(function () {
       ledger.renderGrid()
    })
 }

ledger.exportExcel = function () {
    $("#gridLedger").getKendoGrid().saveAsExcel();
}

ledger.setDate = function () {
    var datepicker = $("#dateStart").data("kendoDatePicker");
    var now = new Date();
    var years = moment(now).format("YYYY")
    var Month = moment(now).format("MM")
    var newDate = years + "-" + Month + "-" + "01"
    datepicker.value(new Date(newDate))
}

ledger.onChangeDateStart = function(val){
    if (val.getTime()>ledger.DateEnd().getTime()){
        ledger.DateEnd(val)
    }
}

ledger.init = function () {
    // ledger.setDate()
    ledger.getAccount()
    ledger.getDataLedger(function () {
        ledger.renderGrid()
        ledger.getDateNow()
    })

}
$(document).ready(function () {
    ledger.init()
})