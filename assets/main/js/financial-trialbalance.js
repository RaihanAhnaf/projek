var tax = {}
tax.beforeBegin = ko.observable();
tax.beforeDebet = ko.observable();
tax.beforeCredit = ko.observable();
tax.beforeTrans = ko.observable();
tax.beforeEnd = ko.observable();
tax.AfterBegin = ko.observable();
tax.AfterTrans = ko.observable();
tax.AfterEnd = ko.observable()
var trial = {}
model.Processing(false)
trial.DateStart = ko.observable(moment().startOf('month').format('YYYY-MM-DD hh:mm'))
trial.DateEnd = ko.observable(new Date)
trial.koFilterIsActive = ko.observable(false);
trial.dataAkiva = ko.observableArray([]);
trial.dataPasiva = ko.observableArray([]);
trial.dataIncome = ko.observableArray([]);
trial.dataTax = ko.observableArray([]);
trial.koGrossBeginingValue = ko.observable();
trial.koGrossTransactiongValue = ko.observable();
trial.koGrossDebetValue = ko.observable();
trial.koGrossCreditValue = ko.observable();
trial.koGrossEndingValue = ko.observable();
trial.valueRevExpBegin = ko.observable();
trial.valueRevExpTrans = ko.observable();
trial.valueRevExpDebet = ko.observable();
trial.valueRevExpCredit = ko.observable();
trial.valueRevExpEnd = ko.observable();
trial.DatePageBar = ko.observable();
trial.dateStart = ko.observable(new Date)
trial.dateEnd = ko.observable(new Date)
trial.filterByValue = ko.observable("Detail")
trial.textSearch = ko.observable()
trial.renderFilterBy = function() {
    var Data = [
        {"value": "Detail", "text": "Detail"},
        {"value": "Period", "text": "Period"}
    ]
    $("#filterBy").kendoDropDownList({
        dataSource: Data,
        filter: "contains",
        dataTextField: "text",
        dataValueField: "value",
        change : function(e) {
            var dataitem = this.dataItem()
            trial.filterByValue(dataitem.value)
        }
    })
}

trial.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    trial.DatePageBar(page)
}

trial.getDataActiva = function (callback) {
    model.Processing(true)
    var url = "/financial/getdataactiva"
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var param = {
        DateStart: dateStart,
        DateEnd: dateEnd,
        Filter: trial.koFilterIsActive()
    }
    ajaxPost(url, param, function (res) {
        trial.dataAkiva(res.Data)
        model.Processing(false)
        callback()
    })
}
trial.renderGridActiva = function () {
    $('#gridActiva').html("")
    var data = trial.dataAkiva()
    var columns = [{
        field: 'ACC_Code',
        title: 'Acc Code',
        width: 50
    }, {
        field: 'Account_Name',
        title: 'Account Name',
        width: 150,
        footerTemplate: "<div style='text-align:center; font-size: 15px;'>Total Aktiva:</div>"
    }, {
        field: 'Begining',
        title: 'Begining',
        template: "#=ChangeToRupiah(Begining)#",
        footerTemplate: "#= ChangeToRupiah(sum) #",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        },
        width: 50
    }, {
        field: 'Debet',
        title: 'Debet',
        template: "#=ChangeToRupiah(Debet)#",
        footerTemplate: "#= ChangeToRupiah(sum) #",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        },
        width: 50
    },{
        field: 'Credit',
        title: 'Credit',
        template: "#=ChangeToRupiah(Credit)#",
        footerTemplate: "#= ChangeToRupiah(sum) #",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        },
        width: 50
    }, {
        field: 'Ending',
        title: 'Ending',
        template: "#=ChangeToRupiah(Ending)#",
        footerTemplate: "#= ChangeToRupiah(sum) #",
        attributes: {
            style: "text-align:right;"
        },
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        width: 50
    }]
    $('#gridActiva').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
                field: "Begining",
                aggregate: "sum"
            }, {
                field: "Debet",
                aggregate: "sum"
            },{
                field: "Credit",
                aggregate: "sum"
            }, {
                field: "Ending",
                aggregate: "sum"
            }],
            sort: [{
                field: "ACC_Code",
                dir: "asc"
            }],
        },
        height: 400,
        scrollable: true,
        columns: columns
    })
}
trial.getDataPasiva = function (callback) {
    model.Processing(true)
    var url = "/financial/getdatapasiva"
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var param = {
        DateStart: dateStart,
        DateEnd: dateEnd,
        Filter: trial.koFilterIsActive()
    }
    ajaxPost(url, param, function (res) {
        trial.dataPasiva(res.Data)
        model.Processing(false)
        callback()
    })
}
trial.rendeGridPasiva = function (refresh, value) {
    $('#gridPasiva').html("")
    var data = trial.dataPasiva()
    // var value4200 = value * -1
    if (refresh) {
        for (i in data) {
            if (data[i].ACC_Code == 4400) {
                if (value>0){
                    data[i].Credit = value
                }else{
                    data[i].Debet = value*-1
                }
                data[i].Transaction = data[i].Debet-data[i].Credit
                data[i].Ending = data[i].Begining + data[i].Credit-data[i].Debet
            }
        }
    }
    var columns = [{
        field: 'ACC_Code',
        title: 'Acc Code',
        width: 50
    }, {
        field: 'Account_Name',
        title: 'Account Name',
        width: 150,
        footerTemplate: "<div style='text-align:center; font-size: 15px;'>Total Pasiva:</div>"
    }, {
        field: 'Begining',
        title: 'Begining',
        template: "#=ChangeToRupiah(Begining)#",
        footerTemplate: "#= ChangeToRupiah(sum) #",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        },
        width: 50
    }, {
        field: 'Debet',
        title: 'Debet',
        template: "#=ChangeToRupiah(Debet)#",
        footerTemplate: "#= ChangeToRupiah(sum) #",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        },
        width: 50
    }, {
        field: 'Credit',
        title: 'Credit',
        template: "#=ChangeToRupiah(Credit)#",
        footerTemplate: "#= ChangeToRupiah(sum) #",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        },
        width: 50
    }, {
        field: 'Ending',
        title: 'Ending',
        template: "#=ChangeToRupiah(Ending)#",
        footerTemplate: "#= ChangeToRupiah(sum) #",
        attributes: {
            style: "text-align:right;"
        },
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        width: 50
    }]
    $('#gridPasiva').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
                field: "Begining",
                aggregate: "sum"
            }, {
                field: "Debet",
                aggregate: "sum"
            }, {
                field: "Credit",
                aggregate: "sum"
            }, {
                field: "Ending",
                aggregate: "sum"
            }],
            sort: [{
                field: "ACC_Code",
                dir: "asc"
            }],
        },
        height: 400,
        scrollable: true,
        columns: columns
    })
}
trial.removeRowGroup = function () {
    $(".k-group-col,.k-group-cell").remove();
    $(".k-grid .k-icon,.k-grid .k-i-collapse").remove();
    var spanCells = $(".k-grouping-row").children("td");
    spanCells.attr("colspan", spanCells.attr("colspan") - 1);
}
trial.getDataIncomeStatement = function (callback) {
    model.Processing(true)
    var url = "/financial/getdatatrialincome"
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var param = {
        DateStart: dateStart,
        DateEnd: dateEnd,
        Filter: trial.koFilterIsActive()
    }
    ajaxPost(url, param, function (res) {
        trial.dataIncome(res.Data)
        model.Processing(false)
        callback()
    })
}
// ================== INCOME ======================
trial.rendergridIncome1 = function () {
    $('#gridIncomeStatement1').html("")
    //SALES AND REVENUE:5000
    var Data5000 = _.filter(trial.dataIncome(), {
        ACC_Code: 5110
    });
    // var Data5600 = _.find(trial.dataIncome(), {
    //     ACC_Code: 5600
    // });
    // if (Data5600 != undefined) {
    //     Data5000.push(Data5600)
    // }
    for (i in Data5000) {
        Data5000[i].Parent = "PENJUALAN"
    }
    var data = Data5000
    var columns = [{
        field: 'ACC_Code',
        title: 'Acc Code',
        width: 50
    }, {
        field: 'Account_Name',
        title: 'Account Name',
        width: 150,
        footerTemplate: "<div style='text-align:left; font-size: 15px;'>TOTAL PENJUALAN:</div>"
    }, {
        field: 'Begining',
        title: 'Begining',
        template: "#=ChangeToRupiah(Begining)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(trial.CalculateGrossBegining()) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Debet',
        title: 'Debet',
        template: "#=ChangeToRupiah(Debet)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Credit',
        title: 'Credit',
        template: "#=ChangeToRupiah(Credit)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Ending',
        title: 'Ending',
        template: "#=ChangeToRupiah(Ending)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #"
    }]
    $('#gridIncomeStatement1').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
                field: "Begining",
                aggregate: "sum"
            }, {
                field: "Debet",
                aggregate: "sum"
            },{
                field: "Credit",
                aggregate: "sum"
            }, {
                field: "Ending",
                aggregate: "sum"
            }],
            group: {
                field: "Parent",
                dir: "asc"
            }
        },
        scrollable: true,
        columns: columns
    })
    var grid = $("#gridIncomeStatement1").data("kendoGrid");
    grid.hideColumn(6);
    $("#gridIncomeSales .k-grid-footer").css('display', 'none');
}
trial.CalculateGrossBegining = function () {
    var data = $("#gridIncomeStatement1").data('kendoGrid').dataSource.options.data
    var total = _.sumBy(data, function (o) {
        return o.Begining;
    });
    trial.koGrossBeginingValue(total)
    return total
}
trial.CalculateGrossTransaction = function () {
    var data = $("#gridIncomeStatement1").data('kendoGrid').dataSource.options.data
    var totalDebet =0.0
    var totalCredit= 0.0
    var totalEnding = 0.0
    for (i in data) {
        totalEnding += data[i].Transaction
        totalDebet +=  data[i].Debet
        totalCredit +=  data[i].Credit
    }
    trial.koGrossDebetValue(totalDebet)
    trial.koGrossCreditValue(totalCredit)
    trial.koGrossTransactiongValue(totalEnding)
    return kendo.toString(totalEnding, 'n')
    // //trans
    // var DataSum = 0.0
    // var DataMin = 0.0
    // //debet
    // var DataSumD = 0.0
    // var DataMinD = 0.0
    // //credit
    // var DataSumC = 0.0
    // var DataMinC = 0.0
    // for (i in data) {
    //     if (data[i].ACC_Code == 5600) {
    //         DataMin += data[i].Transaction
    //         DataMinD +=  data[i].Debet
    //         DataMinC +=  data[i].Credit
    //     } else {
    //         DataSum += data[i].Transaction
    //         DataSumD+=  data[i].Debet
    //         DataSumC+=  data[i].Credit
    //     }
    // }
    // var Resuts = DataSum - DataMin
    // var debet = DataSumD - DataMinD
    // var credit = DataSumC - DataMinC
    // trial.koGrossDebetValue(debet)
    // trial.koGrossCreditValue(credit)
    // trial.koGrossTransactiongValue(Resuts)
    // return kendo.toString(Resuts, 'n')
}
trial.CalculateGrossEnding = function () {
    var data = $("#gridIncomeStatement1").data('kendoGrid').dataSource.options.data
    var DataSum = 0.0
    var DataMin = 0.0
    for (i in data) {
        if (data[i].ACC_Code == 5600) {
            DataMin += data[i].Ending
        } else {
            DataSum += data[i].Ending
        }
    }
    var Resuts = DataSum - DataMin
    trial.koGrossEndingValue(Resuts)
    return kendo.toString(Resuts, 'n')
}
trial.rendergridIncome1Retur = function(){
    $('#gridIncomeStatement1retur').html("")
    var data = _.filter(trial.dataIncome(), function(x){
        return x.ACC_Code== 5210 || x.ACC_Code==5211 || x.ACC_Code==5212  
    });
    for (i in data) {
        data[i].Parent = "HARGA POKOK"
    }
    var columns = [{
        field: 'ACC_Code',
        title: 'Acc Code',
        width: 50
    }, {
        field: 'Account_Name',
        title: 'Account Name',
        width: 150,
        footerTemplate: "<div style='text-align:left; font-size: 15px;'>TOTAL HARGA POKOK:</div><div style='text-align:left; font-size: 15px;'>LABA KOTOR:</div>"
    }, {
        field: 'Begining',
        title: 'Begining',
        template: "#=ChangeToRupiah(Begining)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>-</div>",
        // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= trial.PenjualanBersih('Begining') #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Debet',
        title: 'Debet',
        template: "#=ChangeToRupiah(Debet)#",
        width: 50,
        // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= trial.PenjualanBersih('Debet') #</div>",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>-</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Credit',
        title: 'Credit',
        template: "#=ChangeToRupiah(Credit)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>-</div>",
        // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= trial.PenjualanBersih('Credit') #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Ending',
        title: 'Ending',
        template: "#=ChangeToRupiah(Ending)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= trial.PenjualanBersih('Ending') #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #"
    }]

    $('#gridIncomeStatement1retur').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
                field: "Begining",
                aggregate: "sum"
            }, {
                field: "Debet",
                aggregate: "sum"
            }, {
                field: "Credit",
                aggregate: "sum"
            }, {
                field: "Ending",
                aggregate: "sum"
            }],
            sort: {
                field: "ACC_Code",
                dir: "asc"
            },
            group: {
                field: "Parent",
                dir: "asc"
            }
        },
        scrollable: true,
        columns: columns
    })
    var grid = $("#gridIncomeStatement1retur").data("kendoGrid");
    grid.hideColumn(6);
    $("#gridIncomeStatement1retur > div.k-grid-header > div > table > thead > tr").hide()
}
trial.PenjualanBersih = function(type){
    if (type == "Begining") {
        var totalBegining = $("#gridIncomeStatement1retur").data("kendoGrid").dataSource.aggregates().Begining.sum
        var Resuts = trial.koGrossBeginingValue() - totalBegining
        return ChangeToRupiah(Resuts)
    }
    if (type == "Debet") {
        var totalDebet = $("#gridIncomeStatement1retur").data("kendoGrid").dataSource.aggregates().Debet.sum
        var Resuts = (trial.koGrossDebetValue() - totalDebet)*-1
        return ChangeToRupiah(Resuts)
    }
    if (type == "Credit") {
        var totalCredit = $("#gridIncomeStatement1retur").data("kendoGrid").dataSource.aggregates().Credit.sum
        var Resuts = trial.koGrossCreditValue() - totalCredit
        return ChangeToRupiah(Resuts)
    }
    if (type == "Ending") {
        var totalEnding = $("#gridIncomeStatement1retur").data("kendoGrid").dataSource.aggregates().Ending.sum
        var Resuts = trial.koGrossEndingValue() - totalEnding
        return ChangeToRupiah(Resuts)
    }
}
// trial.rendergridIncome1Hpp = function(){
//     $('#gridIncomeStatement1hpp').html("")
//     var data = _.filter(trial.dataIncome(), function(x){
//         return x.ACC_Code== 5210 
//     });
//     for (i in data) {
//         data[i].Parent = "HARGA POKOK PENJUALAN"
//     }
//     var columns = [{
//         field: 'ACC_Code',
//         title: 'Acc Code',
//         width: 50
//     }, {
//         field: 'Account_Name',
//         title: 'Account Name',
//         width: 150,
//         footerTemplate: "<div style='text-align:left; font-size: 15px;'>TOTAL HARGA POKOK PENJUALAN:</div><div style='text-align:left; font-size: 15px;'>TOTAL LABA KOTOR:</div>"
//     }, {
//         field: 'Begining',
//         title: 'Begining',
//         template: "#=ChangeToRupiah(Begining)#",
//         width: 50,
//         footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>-</div>",
//         // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(trial.labakotor('Begining', sum)) #</div>",
//         attributes: {
//             style: "text-align:right;"
//         }
//     }, {
//         field: 'Debet',
//         title: 'Debet',
//         template: "#=ChangeToRupiah(Debet)#",
//         width: 50,
//         footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>-</div>",
//         // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(trial.labakotor('Debet', sum)) #</div>",
//         attributes: {
//             style: "text-align:right;"
//         }
//     }, {
//         field: 'Credit',
//         title: 'Credit',
//         template: "#=ChangeToRupiah(Credit)#",
//         width: 50,
//         footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>-</div>",
//         // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(trial.labakotor('Credit',sum)) #</div>",
//         attributes: {
//             style: "text-align:right;"
//         }
//     }, {
//         field: 'Ending',
//         title: 'Ending',
//         template: "#=ChangeToRupiah(Ending)#",
//         width: 50,
//         footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(trial.labakotor('Ending', sum)) #</div>",
//         attributes: {
//             style: "text-align:right;"
//         }
//     }, {
//         field: 'Parent',
//         title: 'Parent',
//         width: 50,
//         groupHeaderTemplate: "#= value #"
//     }]

//     $('#gridIncomeStatement1hpp').kendoGrid({
//         dataSource: {
//             data: data,
//             aggregate: [{
//                 field: "Begining",
//                 aggregate: "sum"
//             }, {
//                 field: "Debet",
//                 aggregate: "sum"
//             }, {
//                 field: "Credit",
//                 aggregate: "sum"
//             }, {
//                 field: "Ending",
//                 aggregate: "sum"
//             }],
//             sort: {
//                 field: "ACC_Code",
//                 dir: "asc"
//             },
//             group: {
//                 field: "Parent",
//                 dir: "asc"
//             }
//         },
//         scrollable: true,
//         columns: columns
//     })
//     var grid = $("#gridIncomeStatement1hpp").data("kendoGrid");
//     grid.hideColumn(6);
//     $("#gridIncomeStatement1hpp > div.k-grid-header > div > table > thead > tr").hide()
// }
trial.labakotor = function(type){
    if (type == "Begining") {
        var totalBegining = $("#gridIncomeStatement1retur").data("kendoGrid").dataSource.aggregates().Begining.sum
        // var totalhpp= $("#gridIncomeStatement1hpp").data("kendoGrid").dataSource.aggregates().Begining.sum 
        var penjualanbersih = trial.koGrossBeginingValue() - totalBegining
        // var Results = penjualanbersih - totalhpp
        return penjualanbersih
    }
    if (type == "Debet") {
        var totalDebet = $("#gridIncomeStatement1retur").data("kendoGrid").dataSource.aggregates().Debet.sum
        // var totalhpp= $("#gridIncomeStatement1hpp").data("kendoGrid").dataSource.aggregates().Debet.sum 
        var penjualanbersih = (trial.koGrossDebetValue() - totalDebet) *-1
        // var Results = penjualanbersih - totalhpp
        return penjualanbersih
    }
    if (type == "Credit") {
        var totalCredit = $("#gridIncomeStatement1retur").data("kendoGrid").dataSource.aggregates().Credit.sum
        // var totalhpp= $("#gridIncomeStatement1hpp").data("kendoGrid").dataSource.aggregates().Credit.sum 
        var penjualanbersih = trial.koGrossCreditValue() - totalCredit
        // var Results = penjualanbersih - totalhpp
        return penjualanbersih
    }
    if (type == "Ending") {
        var totalEnding = $("#gridIncomeStatement1retur").data("kendoGrid").dataSource.aggregates().Ending.sum
        // var totalhpp= $("#gridIncomeStatement1hpp").data("kendoGrid").dataSource.aggregates().Ending.sum 
        var penjualanbersih = trial.koGrossEndingValue() - totalEnding
        // var Results = penjualanbersih - totalhpp
        return penjualanbersih
    }
}
trial.rendergridIncome2 = function () {
    $('#gridIncomeStatement2').html("")
    //Operating Expenses:6000
    var Data6000 = _.filter(trial.dataIncome(), function(x){
        return x.ACC_Code!=6160&&x.ACC_Code!=6161&&x.ACC_Code!=6999&& x.Main_Acc_Code == 6000
     });
    for (i in Data6000) {
        Data6000[i].Parent = "BIAYA UMUM DAN ADMINISTRASI"
    }
    var data = Data6000
    var columns = [{
        field: 'ACC_Code',
        title: 'Acc Code',
        width: 50
    }, {
        field: 'Account_Name',
        title: 'Account Name',
        width: 150,
        footerTemplate: "<div style='text-align:left; font-size: 15px;'>TOTAL BIAYA UMUM DAN ADMINISTRASI:</div><div style='text-align:left; font-size: 15px;'>LABA USAHA:</div>"
    }, {
        field: 'Begining',
        title: 'Begining',
        template: "#=ChangeToRupiah(Begining)#",
        width: 50,
        // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(trial.netProfit('Begining')) #</div>",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>-</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Debet',
        title: 'Debet',
        template: "#=ChangeToRupiah(Debet)#",
        width: 50,
        // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(trial.netProfit('Debet')) #</div>",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>-</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Credit',
        title: 'Credit',
        template: "#=ChangeToRupiah(Credit)#",
        width: 50,
        // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(trial.netProfit('Credit')) #</div>",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>-</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Ending',
        title: 'Ending',
        template: "#=ChangeToRupiah(Ending)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(trial.netProfit('Ending')) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #"
    }]

    $('#gridIncomeStatement2').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
                field: "Begining",
                aggregate: "sum"
            }, {
                field: "Debet",
                aggregate: "sum"
            }, {
                field: "Credit",
                aggregate: "sum"
            }, {
                field: "Ending",
                aggregate: "sum"
            }],
            sort: {
                field: "ACC_Code",
                dir: "asc"
            },
            group: {
                field: "Parent",
                dir: "asc"
            }
        },
        scrollable: true,
        columns: columns
    })
    var grid = $("#gridIncomeStatement2").data("kendoGrid");
    grid.hideColumn(6);
    $("#gridIncomeStatement2 > div.k-grid-header > div > table > thead > tr").hide()

}
trial.netProfit = function (type) {
    if (type == "Begining") {
        var totalBegining = $("#gridIncomeStatement2").data("kendoGrid").dataSource.aggregates().Begining.sum
        var labakotor = trial.labakotor(type)
        var Results = labakotor - totalBegining
        return Results
    }
    if (type == "Debet") {
        var totalDebet = $("#gridIncomeStatement2").data("kendoGrid").dataSource.aggregates().Debet.sum
        var labakotor = trial.labakotor(type)
        var Results = labakotor - totalDebet
        return Results
    }
    if (type == "Credit") {
        var totalCredit = $("#gridIncomeStatement2").data("kendoGrid").dataSource.aggregates().Credit.sum
        var labakotor = trial.labakotor(type)
        var Results = labakotor - totalCredit
        return Results
    }
    if (type == "Ending") {
        var totalEnding = $("#gridIncomeStatement2").data("kendoGrid").dataSource.aggregates().Ending.sum
        var labakotor = trial.labakotor(type)
        var Results = labakotor - totalEnding
        return Results
    }
    // if (type == "TransactionFor4400") {
    //     var totalTransaction = $("#gridIncomeStatement2").data("kendoGrid").dataSource.aggregates().Transaction.sum
    //     var Resuts = trial.koGrossBeginingValue() - totalTransaction
    //     return Resuts
    // }
}
trial.rendergridIncome3 = function () {
    $('#gridIncomeStatement3').html("")
    //Pengahasilan lain lain:7000
    var Data7000 = _.filter(trial.dataIncome(), {
        Main_Acc_Code: 7000
    });
    for (i in Data7000) {
        Data7000[i].Parent = "PENDAPATAN/(BIAYA) DI LUAR USAHA"
    }
    var data = Data7000
    var columns = [{
        field: 'ACC_Code',
        title: 'Acc Code',
        width: 50
    }, {
        field: 'Account_Name',
        title: 'Account Name',
        width: 150,
        footerTemplate: "<div style='text-align:left; font-size: 15px;'>TOTAL PENDAPATAN/(BIAYA) DI LUAR USAHA:</div>"
    }, {
        field: 'Begining',
        title: 'Begining',
        template: "#=ChangeToRupiah(Begining)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Debet',
        title: 'Debet',
        template: "#=ChangeToRupiah(Debet)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Credit',
        title: 'Credit',
        template: "#=ChangeToRupiah(Credit)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Ending',
        title: 'Ending',
        template: "#=ChangeToRupiah(Ending)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #"
    }]
    $('#gridIncomeStatement3').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
                field: "Begining",
                aggregate: "sum"
            }, {
                field: "Debet",
                aggregate: "sum"
            }, {
                field: "Credit",
                aggregate: "sum"
            }, {
                field: "Ending",
                aggregate: "sum"
            }],
            group: {
                field: "Parent",
                dir: "asc"
            }
        },
        scrollable: true,
        columns: columns
    })
    var grid = $("#gridIncomeStatement3").data("kendoGrid");
    grid.hideColumn(6);
    $("#gridIncomeStatement3 > div.k-grid-header > div > table > thead > tr").hide()
}
trial.rendergridIncome4 = function () {
    $('#gridIncomeStatement4').html("")
    //Operating Expenses:8000
    var Data8000 = _.filter(trial.dataIncome(), {
        Main_Acc_Code: 8000
    });
    for (i in Data8000) {
        Data8000[i].Parent = "BIAYA DI LUAR USAHA"
    }
    var data = Data8000
    var columns = [{
        field: 'ACC_Code',
        title: 'Acc Code',
        width: 50
    }, {
        field: 'Account_Name',
        title: 'Account Name',
        width: 150,
        footerTemplate: "<div style='text-align:left; font-size: 15px;'>TOTAL BIAYA DI LUAR USAHA:</div><div style='text-align:left; font-size: 15px;'>TOTAL BIAYA DAN PENDAPATAN DI LUAR USAHA:</div><div style='text-align:left; font-size: 15px;'>LABA BERSIH SEBELUM PAJAK:</div>"
    }, {
        field: 'Begining',
        title: 'Begining',
        template: "#=ChangeToRupiah(Begining)#",
        width: 50,
        // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= trial.calculateRevenueExpenses('Begining') #</div><div style='text-align:right; font-size: 15px;'>#= trial.calculateErningBeforTax('Begining') #</div>",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= trial.calculateRevenueExpenses('Begining') #</div><div style='text-align:right; font-size: 15px;'>-</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Debet',
        title: 'Debet',
        template: "#=ChangeToRupiah(Debet)#",
        width: 50,
        // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= trial.calculateRevenueExpenses('Debet') #</div><div style='text-align:right; font-size: 15px;'>#= trial.calculateErningBeforTax('Debet') #</div>",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= trial.calculateRevenueExpenses('Debet') #</div><div style='text-align:right; font-size: 15px;'>-</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Credit',
        title: 'Credit',
        template: "#=ChangeToRupiah(Credit)#",
        width: 50,
        // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= trial.calculateRevenueExpenses('Credit') #</div><div style='text-align:right; font-size: 15px;'>#= trial.calculateErningBeforTax('Credit') #</div>",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= trial.calculateRevenueExpenses('Credit') #</div><div style='text-align:right; font-size: 15px;'>-</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Ending',
        title: 'Ending',
        template: "#=ChangeToRupiah(Ending)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= trial.calculateRevenueExpenses('Ending') #</div><div style='text-align:right; font-size: 15px;'>#= trial.calculateErningBeforTax('Ending') #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #"
    }]
    $('#gridIncomeStatement4').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
                field: "Begining",
                aggregate: "sum"
            }, {
                field: "Debet",
                aggregate: "sum"
            }, {
                field: "Credit",
                aggregate: "sum"
            }, {
                field: "Ending",
                aggregate: "sum"
            }],
            group: {
                field: "Parent",
                dir: "asc"
            }
        },
        scrollable: true,
        columns: columns
    })
    var grid = $("#gridIncomeStatement4").data("kendoGrid");
    grid.hideColumn(6);
    $("#gridIncomeStatement4 > div.k-grid-header > div > table > thead > tr").hide()
}
trial.calculateRevenueExpenses = function (type) {
    if (type == 'Begining') {
        var Revenue = $("#gridIncomeStatement3").data("kendoGrid").dataSource.aggregates().Begining.sum
        var Expenses = $("#gridIncomeStatement4").data("kendoGrid").dataSource.aggregates().Begining.sum
        var Resuts = Revenue - Expenses
        trial.valueRevExpBegin(Resuts)
        return ChangeToRupiah(Resuts)
    }
    if (type == 'Debet') {
        var Revenue = $("#gridIncomeStatement3").data("kendoGrid").dataSource.aggregates().Debet.sum
        var Expenses = $("#gridIncomeStatement4").data("kendoGrid").dataSource.aggregates().Debet.sum
        var Resuts = Revenue - Expenses
        trial.valueRevExpDebet(Resuts)
        return ChangeToRupiah(Resuts)
    }
    if (type == 'Credit') {
        var Revenue = $("#gridIncomeStatement3").data("kendoGrid").dataSource.aggregates().Credit.sum
        var Expenses = $("#gridIncomeStatement4").data("kendoGrid").dataSource.aggregates().Credit.sum
        var Resuts = Revenue - Expenses
        trial.valueRevExpCredit(Resuts)
        return ChangeToRupiah(Resuts)
    }
    if (type == 'Ending') {
        var Revenue = $("#gridIncomeStatement3").data("kendoGrid").dataSource.aggregates().Ending.sum
        var Expenses = $("#gridIncomeStatement4").data("kendoGrid").dataSource.aggregates().Ending.sum
        var Resuts = Revenue - Expenses
        trial.valueRevExpEnd(Resuts)
        return ChangeToRupiah(Resuts)
    }
}
trial.calculateErningBeforTax = function (type) {
    if (type == "Begining") {
        // var totalBegining = $("#gridIncomeStatement2").data("kendoGrid").dataSource.aggregates().Begining.sum
        var netProfit = trial.netProfit(type)
        var Resuts = netProfit + trial.valueRevExpBegin()
        tax.beforeBegin(Resuts)
        return ChangeToRupiah(Resuts)
    }
    if (type == "Debet") {
        // var totalTransaction = $("#gridIncomeStatement2").data("kendoGrid").dataSource.aggregates().Debet.sum
        var netProfit = trial.netProfit(type)
        var Resuts = netProfit + trial.valueRevExpDebet()
        tax.beforeDebet(Resuts)
        return ChangeToRupiah(Resuts)
    }
    if (type == "Credit") {
        // var totalTransaction = $("#gridIncomeStatement2").data("kendoGrid").dataSource.aggregates().Credit.sum
        var netProfit = trial.netProfit(type)
        var Resuts = netProfit + trial.valueRevExpCredit()
        tax.beforeCredit(Resuts)
        return ChangeToRupiah(Resuts)
    }
    if (type == "Ending") {
        // var totalEnding = $("#gridIncomeStatement2").data("kendoGrid").dataSource.aggregates().Ending.sum
        var netProfit = trial.netProfit(type)
        var Resuts = netProfit + trial.valueRevExpEnd()
        tax.beforeEnd(Resuts)
        return ChangeToRupiah(Resuts)
    }
}
trial.EndingCalculation = function (id) {

    var data = $("#" + id).data('kendoGrid').dataSource.options.data
    var sumDebit = _.sumBy(data, 'Begining')
    var sumCredt = _.sumBy(data, 'Transaction')
    var Saldo = sumDebit + sumCredt
    return ChangeToRupiah(Saldo)
}
trial.goIncomegrid = function (grid) {
    if (grid == 1) {
        $('html, body').animate({
            scrollTop: $("#gridIncomeStatement1").offset().top
        }, 700);
    }
    if (grid == 2) {
        $('html, body').animate({
            scrollTop: $("#gridIncomeStatement2").offset().top
        }, 700);
    }
    if (grid == 3) {
        $('html, body').animate({
            scrollTop: $("#gridIncomeStatement3").offset().top
        }, 1000);
    }
    if (grid == 4) {
        $('html, body').animate({
            scrollTop: $("#gridIncomeStatement4").offset().top
        }, 1000);
    }

}
trial.getDataTax = function (callback) {
    model.Processing(true)
    var url = "/financial/getdatatax"
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var param = {
        DateStart: dateStart,
        DateEnd: dateEnd,
        Filter: trial.koFilterIsActive()
    }
    ajaxPost(url, param, function (res) {
        trial.dataTax(res.Data)
        model.Processing(false)
        callback()
    })
}
trial.renderGridTAX = function () {
    $('#gridTAX').html("")
    var dataTax = _.filter(trial.dataTax(), {
        Account_Name: "PAJAK"
    });
    var data = []
    if (dataTax.length == 0) {
        data = [{
            "ACC_Code": "",
            "Account_Name": "TAX",
            "Begining": 0,
            "Category": "TAX",
            "Debet_Credit": "TAX",
            "Ending": 0,
            "ID": "5925503a92de4d5928d773d3",
            "Main_Acc_Code": 2000,
            "Transaction": 0,
            "Parent": "TAX"
        }]
    } else {
        data = dataTax
    }
    var columns = [{
        field: 'ACC_Code',
        title: 'Acc Code',
        width: 50
    }, {
        field: 'Account_Name',
        title: 'Account Name',
        width: 150,
        footerTemplate: "<div style='text-align:left; font-size: 15px;'>EARNING AFTER TAX:</div>"
    }, {
        field: 'Begining',
        title: 'Begining',
        template: "#=ChangeToRupiah(Begining)#",
        width: 50,
        // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= trial.calculateErningAfterTax('begining',sum) #</div>",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>-</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Debet',
        title: 'Debet',
        template: "#=ChangeToRupiah(Debet)#",
        width: 50,
        // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= trial.calculateErningAfterTax('Debet',sum) #</div>",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>-</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Credit',
        title: 'Credit',
        template: "#=ChangeToRupiah(Credit)#",
        width: 50,
        // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= trial.calculateErningAfterTax('Credit',sum) #</div>",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>-</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Ending',
        title: 'Ending',
        template: "#=ChangeToRupiah(Ending)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= trial.calculateErningAfterTax('ending',sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, ]
    $('#gridTAX').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
                field: "Begining",
                aggregate: "sum"
            }, {
                field: "Debet",
                aggregate: "sum"
            }, {
                field: "Credit",
                aggregate: "sum"
            }, {
                field: "Ending",
                aggregate: "sum"
            }],
        },
        scrollable: true,
        columns: columns
    })
    $("#gridTAX > div.k-grid-header > div > table > thead > tr").hide()
    // var valueFor4400 = trial.netProfit('TransactionFor4400')

}
trial.calculateErningAfterTax = function (type, value) {
    if (type == "begining") {
        var Resuts = tax.beforeBegin() - value
        return ChangeToRupiah(Resuts)
    }
    if (type == "Debet") {
        var Resuts = tax.beforeDebet() - value
        return ChangeToRupiah(Resuts)
    }
    if (type == "Credit") {
        var totaldebet = $("#gridTAX").data("kendoGrid").dataSource.aggregates().Debet.sum
        var resultDebet = tax.beforeDebet()- totaldebet
        var Resuts = tax.beforeCredit() - value
        return ChangeToRupiah(Resuts)
    }
    if (type == "ending") {
        var Resuts = tax.beforeEnd() - value
        trial.rendeGridPasiva(true,Resuts)
        return ChangeToRupiah(Resuts)
    }
}
trial.setDate = function () {
    var datepicker = $("#dateStart").data("kendoDatePicker");
    var now = new Date();
    var years = moment(now).format("YYYY")
    var Month = moment(now).format("MM")
    var newDate = years + "-" + Month + "-" + "01"
    datepicker.value(new Date(newDate))
    trial.dateStart(new Date(newDate))
}
trial.showColumn5BeforeRefresh = function () {
    var grid1 = $("#gridIncomeStatement1").data("kendoGrid");
    grid1.showColumn(5);
    var grid2 = $("#gridIncomeStatement2").data("kendoGrid");
    grid2.showColumn(5);
    var grid3 = $("#gridIncomeStatement3").data("kendoGrid");
    grid3.showColumn(5);
    var grid4 = $("#gridIncomeStatement4").data("kendoGrid");
    grid4.showColumn(5);
    $("#gridTAX > div.k-grid-header > div > table > thead > tr").show()
}
trial.refreshGrid = function (e) {
    trial.textSearch(e)
    trial.showColumn5BeforeRefresh()
    trial.koFilterIsActive(true)
    if(trial.filterByValue() == "Detail") {
        trial.GetDataEarning()
        trial.getDataActiva(function () {
            trial.renderGridActiva()
        })
        trial.getDataPasiva(function () {
            trial.rendeGridPasiva(false,0)
        })
        trial.getDataIncomeStatement(function () {
            trial.rendergridIncome1()
            trial.CalculateGrossBegining()
            trial.CalculateGrossTransaction()
            trial.CalculateGrossEnding()
            trial.rendergridIncome1Retur()
            // trial.rendergridIncome1Hpp()
            trial.rendergridIncome2()
            trial.rendergridIncome3()
            trial.rendergridIncome4()

            trial.getDataTax(function () {
                trial.renderGridTAX()
            })
            trial.removeRowGroup()
        })
    }else {
        // console.log("asd")
        trial2.getDataActivaPeriod(function () {
            trial2.renderGridActivaPeriod()
        })
        trial2.getDataPasivaPeriod(function() {
            trial2.getRenderPasivaPeriod()
        })
        trial2.getDataIncomePeriod(function () {
            trial2.rendergridIncome1()
            trial2.rendergridIncome1Retur()
            trial2.rendergridIncome2()
            trial2.rendergridIncome3()
            trial2.rendergridIncome4()
        })
        trial2.getDataTaxPeriode(function () {
            trial2.renderGridTax()
        })
        trial.removeRowGroup()
        // trial2.getDataPeriod = function() {
        //     trial2.renderGridActivaPeriod()
        // }
    }

}
trial.GetDataEarning = function () {
    var url = "/financial/getincomestatementforearning"
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var param = {
        DateStart: dateStart,
        DateEnd: dateEnd,
    }
    ajaxPost(url, param, function (res) {})
}

trial.onChangeDateStart = function(val){
    if (val.getTime()>trial.dateEnd().getTime()){
        trial.dateEnd(val)
    }
}
trial.init = function () {
    trial.getDateNow()
    trial.setDate()
    trial.renderFilterBy()
    // trial.GetDataEarning()
    trial.getDataActiva(function () {
        trial.renderGridActiva()
    })
    trial.getDataPasiva(function () {
        trial.rendeGridPasiva(false,0)
    })
    trial.getDataIncomeStatement(function () {
        trial.rendergridIncome1()
        trial.CalculateGrossBegining()
        trial.CalculateGrossTransaction()
        trial.CalculateGrossEnding()
        trial.rendergridIncome1Retur()
        // trial.rendergridIncome1Hpp()
        trial.rendergridIncome2()
        trial.rendergridIncome3()
        trial.rendergridIncome4()

        trial.getDataTax(function () {
            trial.renderGridTAX()
        })
        trial.removeRowGroup()
    })
}
$(document).ready(function () {
    trial.init()
})