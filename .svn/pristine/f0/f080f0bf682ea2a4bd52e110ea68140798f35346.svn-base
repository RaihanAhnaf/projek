var income = {}
income.DateStart = ko.observable(moment().startOf('month').format('YYYY-MM-DD hh:mm'))
income.DateEnd = ko.observable(new Date)
income.dataMasterIncome = ko.observableArray([]);
income.dataTax = ko.observableArray([]);
income.koFilterIsActive = ko.observable(false);
income.textSearch = ko.observable();
income.koGrossEndingValue = ko.observable(0);
income.valueRevExpEnd = ko.observable(0);
income.taxbeforeEnd = ko.observable();
income.salesEnding = ko.observable(0);
income.totalPotonganAndRetur = ko.observable(0); 
income.netProfitPercentage = ko.observable(0);
income.calculateRevenueExpensesPercentage = ko.observable(0);
income.calculateEarningBeforTaxPercentage = ko.observable(0);
income.calculateEarningAfterTaxPercentage = ko.observable(0);
income.totalSalesAndRevenue = ko.observable(0)
income.TitelFilter = ko.observable(" Hide Filter")
income.DatePageBar = ko.observable()
income.valueTotalPenjualanBersih = ko.observable(0)
income.valueTotalHargaPokokPenjualan = ko.observable(0)
income.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    income.DatePageBar(page)
}

income.removeRowGroup = function () {
    $(".k-group-col,.k-group-cell").remove();
    $(".k-grid .k-icon,.k-grid .k-i-collapse").remove();
    var spanCells = $(".k-grouping-row").children("td");
    spanCells.attr("colspan", spanCells.attr("colspan") - 1);
}

income.getDataTax = function (callback) {
    model.Processing(true)
    var url = "/financial/getdatatax"
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var param = {
        DateStart: dateStart,
        DateEnd: dateEnd,
        Filter: income.koFilterIsActive()
    }
    ajaxPost(url, param, function (res) {
        income.dataTax(res.Data)
        model.Processing(false)
        callback()
    })
}

income.getData = function (callback) {
    model.Processing(true)
    var url = "/financial/getdatatrialincome"
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var param = {}
    var v = income.dataMasterIncome()
    if (income.koFilterIsActive() == true) {
        param = {
            DateStart: dateStart,
            DateEnd: dateEnd,
            TextSearch : income.textSearch(),
            Filter: true,            
        }
    }else{
        param = {
            Filter: false,
        }
    }
    ajaxPost(url, param, function (res) {
        model.Processing(false)
        if (res.IsError === "true") {
            swal("Error!", res.Message, "error")
            return
        }

        income.dataMasterIncome(res.Data)

        callback()
    }, function () {
        swal("Error!", "Unknown error, please try again", "error")
    })
}

income.rendergridSales = function () {
    //SALES AND REVENUE:5000
    var Data5000 = _.filter(income.dataMasterIncome(), function(x){
        return x.ACC_Code== 5110
    });
    // var Data5100 = _.filter(income.dataMasterIncome(), {
    //     ACC_Code: 5100
    // });
    // console.log(Data5100)

    income.salesEnding(Data5000[0].Ending)

    // var Data5000 = _.filter(income.dataMasterIncome(), {
    //     Main_Acc_Code: 5000
    // });
    // var Data5600 = _.find(income.dataMasterIncome(), {
    //     ACC_Code: 5600
    // });

    // if (Data5600 != undefined) {
    //     Data5000.push(Data5600)
    // }
    for (i in Data5000) {
        Data5000[i].Parent = "SALES AND REVENUE"
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
        field: 'Ending',
        title: 'Ending',
        template: "#=ChangeToRupiah(Ending)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: '',
        title: 'Sales (%)',
        template: "#= income.salesValue(Ending)#%",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#=income.sumSalesPercentage()#%</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #",
        hidden: true,
    }]

    $('#gridIncomeSales').html("")
    $('#gridIncomeSales').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
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
    // $("#gridIncomeSales .k-grid-footer").css('display', 'none');
}
income.salesValue = function (e) {
    // var Data5000 = _.filter(income.dataMasterIncome(), {
    //     Main_Acc_Code: 5000
    // });
    // var Data5600 = _.find(income.dataMasterIncome(), {
    //     ACC_Code: 5600
    // });

    // if (Data5600 != undefined) {
    //     Data5000.push(Data5600)
    // }
    var Data5000 = _.filter(income.dataMasterIncome(), function(x){
        return x.ACC_Code== 5110
    });
    for (i in Data5000) {
        Data5000[i].Parent = "PENJUALAN"
    }
    var data = Data5000

    var total = _.sumBy(data, function (o) {
        return o.Ending;
    });
    income.totalSalesAndRevenue(total)
    if (e == 0) {
        return 0
    } else {
        var fixValue = div(e, total)
        var salesValues = (fixValue * 100)
        if (salesValues == 0) {
            return salesValues
        }
        return salesValues.toFixed(2)

    }
    return 0
}

income.sumSalesPercentage = function () {
    var totalPercentage = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
    income.koGrossEndingValue(totalPercentage)
    if (totalPercentage > 0) {  
        var Results = (totalPercentage / totalPercentage) * 100
        return ChangeToRupiah(Results)
    }

    return 0
}
income.rendergridIncomeRetur = function(){
    var data = _.filter(income.dataMasterIncome(), function(x){
        return x.ACC_Code == 5210 || x.ACC_Code==5211 || x.ACC_Code==5212
    });
    var totalIncomeRetur = _.sumBy(data, function(o) { return o.Ending; })
    income.totalPotonganAndRetur(totalIncomeRetur)
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
        field: 'Ending',
        title: 'Ending',
        template: "#=ChangeToRupiah(Ending)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income.totalpenjualanbersih(sum)) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: '',
        title: 'Sales (%)',
        template: "#= income.salesReture(Ending)#%",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income.sumSalesReture())#%</div><div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income.totalpercentagepenjualanbersih()) #%</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #",
        hidden: true,
    }]

    $('#gridIncomeRetur').html("")
    $('#gridIncomeRetur').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
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
    $("#gridIncomeRetur > div.k-grid-header > div > table > thead > tr").hide()
}
income.salesReture = function(e){
    var totalPercentage = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
    if (totalPercentage > 0) {  
        var Results = (e / totalPercentage) * 100
        return ChangeToRupiah(Results)
    }

    return 0
}
income.sumSalesReture = function(){
    var totalPercentagePenjualan = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
    var totalPercentageSalesRetur = $("#gridIncomeRetur").data("kendoGrid").dataSource.aggregates().Ending.sum
    if (totalPercentagePenjualan > 0) {  
        var Results = (totalPercentageSalesRetur / totalPercentagePenjualan) * 100
        return Results
    }
    return 0
}
income.totalpenjualanbersih = function(sum){
    var totalPenjualan = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
    var results = totalPenjualan- sum
    income.valueTotalPenjualanBersih(results)
    return results
}
income.totalpercentagepenjualanbersih = function(){
    var totalPenjualan = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
    var sumthisgrid = $("#gridIncomeRetur").data("kendoGrid").dataSource.aggregates().Ending.sum
    if (totalPenjualan > 0) { 
        var res = totalPenjualan- sumthisgrid 
        var Results = (res / totalPenjualan) * 100
        return Results
    }
    return 0
}
// income.rendergridIncomeStatement1_5 = function () {
//     //SALES AND REVENUE:5000
//     var Data5000 = _.filter(income.dataMasterIncome(), function(x){
//         return x.ACC_Code== 5210
//     });
// // 
//     income.salesEnding(_.sumBy(Data5000, function(o) { return o.Ending; }))
//     for (i in Data5000) {
//         Data5000[i].Parent = "HARGA POKOK PENJUALAN"
//     }
//     var data = Data5000
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
//         field: 'Ending',
//         title: 'Ending',
//         template: "#=ChangeToRupiah(Ending)#",
//         width: 50,
//         footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income.totalLabaKotor(sum)) #</div>",
//         attributes: {
//             style: "text-align:right;"
//         }
//     }, {
//         field: '',
//         title: 'Sales (%)',
//         template: "#= income.salesHPP(Ending)#%",
//         width: 50,
//         footerTemplate: "<div style='text-align:right; font-size: 15px;'>#=income.totalPercentageHargaPokokPenjualan()#%</div><div style='text-align:right; font-size: 15px;'>#=income.totalPercentageLabaKotor()#%</div>",
//         attributes: {
//             style: "text-align:right;"
//         }
//     }, {
//         field: 'Parent',
//         title: 'Parent',
//         width: 50,
//         groupHeaderTemplate: "#= value #",
//         hidden: true,
//     }]

//     $('#gridIncomeStatement1_5').html("")
//     $('#gridIncomeStatement1_5').kendoGrid({
//         dataSource: {
//             data: data,
//             aggregate: [{
//                 field: "Ending",
//                 aggregate: "sum"
//             }],
//             group: {
//                 field: "Parent",
//                 dir: "asc"
//             }
//         },
//         scrollable: true,
//         columns: columns
//     })
//     $("#gridIncomeStatement1_5 > div.k-grid-header > div > table > thead > tr").hide()
// }
income.salesHPP = function(e){
    var totalPenjualan = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
    if (totalPenjualan > 0) {  
        var Results = (e / totalPenjualan) * 100
        return ChangeToRupiah(Results)
    }

    return 0
}
// income.totalHargaPokokPenjualan = function(sum){
//     var results = income.valueTotalPenjualanBersih() - sum
//     income.valueTotalHargaPokokPenjualan( results)
//     return ChangeToRupiah(results) 
// }
income.totalPercentageHargaPokokPenjualan = function(){
    var totalPenjualan = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
    var sum =  $("#gridIncomeRetur").data("kendoGrid").dataSource.aggregates().Ending.sum
    // var results = income.valueTotalPenjualanBersih() - sum
    if (totalPenjualan > 0) { 
        var Results = (sum / totalPenjualan) * 100
        return ChangeToRupiah(Results)
    }
    return 0
}
income.totalLabaKotor = function(sum){
    var sumgridpenjualanbersih = $("#gridIncomeRetur").data("kendoGrid").dataSource.aggregates().Ending.sum
    var totalpenjualanbersih = income.totalpenjualanbersih(sumgridpenjualanbersih)
    // var totalHargaPokokPenjualan = income.valueTotalHargaPokokPenjualan()
    var result = totalpenjualanbersih - sum
    return result
}
income.totalPercentageLabaKotor = function(){
    var totalPenjualan = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
    // var totalPotonganDanRetur = $("#gridIncomeRetur").data("kendoGrid").dataSource.aggregates().Ending.sum
    // var totalHargaPokokPenjualan = income.valueTotalHargaPokokPenjualan()
    var sum =  $("#gridIncomeRetur").data("kendoGrid").dataSource.aggregates().Ending.sum
    var totalLabaKotor =income.totalLabaKotor(sum)
    if (totalPenjualan > 0) { 
        var Results = (totalLabaKotor / totalPenjualan) * 100
        return ChangeToRupiah(Results)
    }
    return 0
}
// income.sumSalesPercentageLabaKotor = function(){
//     var totalPercentage = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
//     if (income.koGrossEndingValue() > 0) {  
//         var Results = (totalPercentage /  income.koGrossEndingValue()) * 100
//         return ChangeToRupiah(Results.toFixed(2))
//     }
//     return 0
// }
// income.sumSalesEndingLabaKotor = function(){
    

//     // var v5210 = _.find(income.dataMasterIncome(), function(x){
//     //         return x.ACC_Code== 5210
//     // });
//     var v5210 = $("#gridIncomeStatement1_5").data("kendoGrid").dataSource.aggregates()
//     var totalEnding = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
//     // income.koGrossEndingValue(totalEnding - v5210.Ending)
//     return totalEnding - v5210.Ending.sum
// }

// income.calculateGrossEnding = function () {
//     var data = $("#gridIncomeSales").data('kendoGrid').dataSource.options.data
//     var DataSum = 0.0
//     var DataMin = 0.0
//     for (i in data) {
//         // if (data[i].ACC_Code == 5600) {
//         //     DataMin += data[i].Ending
//         // } else {
//             DataSum += data[i].Ending
//         // }
//     }
//     var Results = DataSum 
//     income.koGrossEndingValue(Results)
//     return kendo.toString(Results, 'n')
// }



income.rendergridOperatingExpenses = function () {
    //Operating Expenses:6000
    // var Data6000 = _.filter(income.dataMasterIncome(), {
    //     Main_Acc_Code: 6000
    // });
    var Data6000 = _.filter(income.dataMasterIncome(), function(x){
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
        field: 'Ending',
        title: 'Ending',
        template: "#=ChangeToRupiah(Ending)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income.totalLabaUsaha(sum)) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: '',
        title: 'Sales (%)',
        template: "#= income.salesOpex(Ending)#%",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income.sumOperatingxpensesPercentage())#%</div><div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income.totalSalesLabaUsaha())#%</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #",
        hidden: true,
    }]
    $('#gridIncomeOperatingExpenses').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
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

    $("#gridIncomeOperatingExpenses > div.k-grid-header > div > table > thead > tr").hide()
}
income.salesOpex = function(e){
    var totalPenjualan = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
    if (totalPenjualan > 0) {  
        var Results = (e / totalPenjualan) * 100
        return ChangeToRupiah(Results)
    }

    return 0
}
income.sumOperatingxpensesPercentage = function () {
    var totalPenjualan = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
    if (totalPenjualan > 0) {
        var totalPercentage = $("#gridIncomeOperatingExpenses").data("kendoGrid").dataSource.aggregates().Ending.sum
        var Results = (totalPercentage / totalPenjualan) * 100
        return Results
    }

    return 0
}
income.totalLabaUsaha = function(sum){
    var sumIncomeStatement1_5 =  $("#gridIncomeRetur").data("kendoGrid").dataSource.aggregates().Ending.sum
    var labakotor = income.totalpenjualanbersih(sumIncomeStatement1_5)
    var result =labakotor- sum
    return result
} 
income.totalSalesLabaUsaha = function(){
    var totalPenjualan = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
    //laba usaha 
    var totalPotonganDanRetur = $("#gridIncomeRetur").data("kendoGrid").dataSource.aggregates().Ending.sum
    var sum = $("#gridIncomeOperatingExpenses").data("kendoGrid").dataSource.aggregates().Ending.sum
    var totalHargaPokokPenjualan = income.valueTotalHargaPokokPenjualan()
    var labakotor = totalPotonganDanRetur - totalHargaPokokPenjualan
    var result =labakotor- sum
    //
    if (totalPenjualan > 0) {
        var totalPercentage = $("#gridIncomeOperatingExpenses").data("kendoGrid").dataSource.aggregates().Ending.sum
        var Results = (result / totalPenjualan) * 100
        return Results
    }

    return 0
}
// income.netProfit = function (type) {
//     var v5210 = $("#gridIncomeStatement1_5").data("kendoGrid").dataSource.aggregates().Ending.sum
//     if (type == "Ending") {
//         var totalEnding = $("#gridIncomeOperatingExpenses").data("kendoGrid").dataSource.aggregates().Ending.sum
//         var labakotor = income.koGrossEndingValue() -v5210
//         var Results = labakotor - totalEnding
//         // var totalEnding = $("#gridIncomeStatement1").data("kendoGrid").dataSource.aggregates().Ending.sum
//         // trial.koGrossEndingValue(totalEnding - v5210.Ending.sum)
//         if (income.koGrossEndingValue() > 0) {
//             var g = (Results / income.koGrossEndingValue()) * 100
//             income.netProfitPercentage(g.toFixed(2))
//         } else {
//             income.netProfitPercentage(0)
//         }
//         return ChangeToRupiah(Results)
//     }
//     if (type == "TransactionFor4400") {
//         var totalTransaction = $("#gridIncomeOperatingExpenses").data("kendoGrid").dataSource.aggregates().Transaction.sum
//         var Results = income.koGrossBeginingValue() - totalTransaction
//         return Results
//     }
// }

income.rendergridOtherIncome = function () {
    //Pengahasilan lain lain:7000
    // var Data7000 = _.filter(income.dataMasterIncome(), {
    //     Main_Acc_Code: 7000
    // });
    var Data7000 = _.filter(income.dataMasterIncome(),function(x){
        return x.ACC_Code == 7100 ||x.ACC_Code == 7200||x.ACC_Code == 7999
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
        footerTemplate: "<div style='text-align:left; font-size: 15px;'>TOTAL PENDAPATAN DI LUAR USAHA:</div>"
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
        field: '',
        title: 'Sales (%)',
        template: "#= income.salesOtherIncome(Ending)#%",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= income.sumOtherIncomePercentage()#%</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #",
        hidden: true
    }]
    $('#gridIncomeOtherIncome').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
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
    $("#gridIncomeOtherIncome > div.k-grid-header > div > table > thead > tr").hide()
}
income.salesOtherIncome = function(e){
    var totalPenjualan = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
    if (totalPenjualan > 0) {  
        var Results = (e / totalPenjualan) * 100
        return ChangeToRupiah(Results)
    }

    return 0
}
income.sumOtherIncomePercentage = function () {
    var totalPenjualan = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
    if (totalPenjualan > 0) {
        var totalPercentage = $("#gridIncomeOtherIncome").data("kendoGrid").dataSource.aggregates().Ending.sum
        var Results = (totalPercentage / totalPenjualan) * 100
        return ChangeToRupiah(Results)
    }

    return 0
}

income.rendergridOtherExpense = function () {
    //Operating Expenses:8000
    // var Data8000 = _.filter(income.dataMasterIncome(), {
    //     Main_Acc_Code: 8000
    // });
    var Data8000 = _.filter(income.dataMasterIncome(), function(x){
       return x.ACC_Code==8100||x.ACC_Code==8200||x.ACC_Code==8300||x.ACC_Code==8400|| x.ACC_Code == 8500
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
        field: 'Ending',
        title: 'Ending',
        template: "#=ChangeToRupiah(Ending)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right; font-size: 15px;'>#= income.calculateRevenueExpenses('ending') #</div><div style='text-align:right; font-size: 15px;'>#= income.calculateEarningBeforTax('ending') #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: '',
        title: 'Sales (%)',
        template: "#= income.salesIncomeOtherExpenses(Ending)#%",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= income.sumOtherExpensesPercentage()#%</div><div style='text-align:right; font-size: 15px;'>#= income.calculateRevenueExpensesPercentage()#%</div><div style='text-align:right; font-size: 15px;'>#= income.calculateEarningBeforTaxPercentage()#%</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #",
        hidden: true
    }]
    $('#gridIncomeOtherExpenses').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
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

    $("#gridIncomeOtherExpenses > div.k-grid-header > div > table > thead > tr").hide()
}
income.salesIncomeOtherExpenses = function(e){
    var totalPenjualan = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
    if (totalPenjualan > 0) {  
        var Results = (e / totalPenjualan) * 100
        return ChangeToRupiah(Results)
    }

    return 0
}
income.sumOtherExpensesPercentage = function () {
    var totalPenjualan = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
    if (totalPenjualan > 0) {
        var totalPercentage = $("#gridIncomeOtherExpenses").data("kendoGrid").dataSource.aggregates().Ending.sum
        var Results = (totalPercentage / totalPenjualan) * 100
        return ChangeToRupiah(Results)
    }

    return 0
}

income.calculateRevenueExpenses = function (type) {
    var totalPenjualan = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
    if (type == 'ending') {
        var Revenue = $("#gridIncomeOtherIncome").data("kendoGrid").dataSource.aggregates().Ending.sum
        var Expenses = $("#gridIncomeOtherExpenses").data("kendoGrid").dataSource.aggregates().Ending.sum
        var Results = Revenue - Expenses
        income.valueRevExpEnd(Results)
        if (totalPenjualan > 0) {
            var g = (Results / totalPenjualan) * 100
            // income.calculateEarningBeforTaxPercentage(g.toFixed(2))
            income.calculateRevenueExpensesPercentage(g.toFixed(2))
            if(g <0){
                // income.calculateEarningBeforTaxPercentage(g.toFixed(2))
                income.calculateRevenueExpensesPercentage(ChangeToRupiah(g.toFixed(2)))
            }
        } else {
            income.calculateRevenueExpensesPercentage(0)
        }
        return ChangeToRupiah(Results)
    }
}

income.calculateEarningBeforTax = function (type) {
    var totalPenjualan = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
    if (type == "ending") {
        // var totalPenjualan = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
        // var totalEnding = $("#gridIncomeOperatingExpenses").data("kendoGrid").dataSource.aggregates().Ending.sum
        // var v5210 = $("#gridIncomeStatement1_5").data("kendoGrid").dataSource.aggregates().Ending.sum
        // var labakotor = totalPenjualan -v5210
        // var netProfit = labakotor - totalEnding
        // // var netProfit = income.koGrossEndingValue() - totalEnding
        // var Results = netProfit + income.valueRevExpEnd()
        // income.taxbeforeEnd(Results)
        // if (income.koGrossEndingValue() > 0) {
        //     var g = (Results / income.koGrossEndingValue()) * 100
        //     income.calculateEarningBeforTaxPercentage(g.toFixed(2))

        // } else {
        //     income.calculateEarningBeforTaxPercentage(0)
        // }
        // return ChangeToRupiah(Results)
        var Revenue = $("#gridIncomeOtherIncome").data("kendoGrid").dataSource.aggregates().Ending.sum
        var Expenses = $("#gridIncomeOtherExpenses").data("kendoGrid").dataSource.aggregates().Ending.sum
        var Opex = $("#gridIncomeOperatingExpenses").data("kendoGrid").dataSource.aggregates().Ending.sum
        var labausaha = income.totalLabaUsaha(Opex)
        var ResultsRevEx = Revenue - Expenses
        var Results = labausaha +ResultsRevEx
        income.taxbeforeEnd(Results)
        if (totalPenjualan > 0) {
            var g = (Results / totalPenjualan) * 100
            income.calculateEarningBeforTaxPercentage(g.toFixed(2))
            if(g <0){
                income.calculateEarningBeforTaxPercentage(ChangeToRupiah(g.toFixed(2)))
            }

        } else {
            income.calculateEarningBeforTaxPercentage(0)
        }
        return ChangeToRupiah(Results)
    }
}

income.renderGridTAX = function () {
    var dataTax = _.filter(income.dataTax(), {
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
        footerTemplate: "<div style='text-align:left; font-size: 15px;'>LABA BERSIH SETALAH PAJAK:</div>"
    }, {
        field: 'Ending',
        title: 'Ending',
        template: "#=ChangeToRupiah(Ending)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= income.calculateEarningAfterTax('ending',sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: '',
        title: 'Sales (%)',
        template: "#= income.salesValue(Ending)#%",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= income.calculateEarningAfterTaxPercentage()#%</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #",
        hidden: true
    }

    ]
    $('#gridTAX').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
                field: "Ending",
                aggregate: "sum"
            }],
        },
        scrollable: true,
        columns: columns
    })
    $("#gridTAX > div.k-grid-header > div > table > thead > tr").hide()

}

income.calculateEarningAfterTax = function (type, value) {
    var totalPenjualan = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates().Ending.sum
    if (type == "ending") {
        var Results = income.taxbeforeEnd() - value
        if (totalPenjualan > 0) {
            var g = (Results / income.totalSalesAndRevenue()) * 100
            income.calculateEarningAfterTaxPercentage(g.toFixed(2))
            if(g <0){
                income.calculateEarningAfterTaxPercentage(ChangeToRupiah(g.toFixed(2)))
            }
        } else {
            income.calculateEarningAfterTaxPercentage(0)
        }
        return ChangeToRupiah(Results)
    }
}

income.setDate = function () {
    var datepicker = $("#dateStart").data("kendoDatePicker");
    var now = new Date();
    var years = moment(now).format("YYYY")
    var Month = moment(now).format("MM")
    var newDate = years + "-" + Month + "-" + "01"
    datepicker.value(new Date(newDate))
}
income.ExportToPDF = function () {
    model.Processing(true)
    var url = "/financial/exporttopdfincome"
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var param = {
        DateStart: dateStart,
        DateEnd: dateEnd
    }
    ajaxPost(url, param, function (res) {
        model.Processing(false)
        window.open('/res/docs/incomestatement/' + res, '_blank');
    })
}


income.refreshGrid = function (e) {
    income.textSearch(e)
    income.koFilterIsActive(true)
    if(income2.filterByValue()== "Detail"){
        income.getData(function () {
            income.rendergridSales()
            income.rendergridIncomeRetur()
            // income.rendergridIncomeStatement1_5()
            income.rendergridOperatingExpenses()
            income.rendergridOtherIncome()
            income.rendergridOtherExpense()
            income.removeRowGroup()
            income.getDataTax(function () {
                income.renderGridTAX()
            })
        })
        income2.buttonExportExcel(false)
    }else{
        income2.getDataPeriode(function () {
            income2.rendergridSalesPeriod()
            income2.rendergridIncomeRetur()
            income2.rendergridOperatingExpensesPeriod()
            income2.rendergridOtherIncomePeriod()
            income2.rendergridOtherExpensePeriod()
            income.removeRowGroup()
            income2.getDataTaxPeriod(function () {
                income2.renderGridTAXPeriod()
            })
        })
        income2.buttonExportExcel(true)
    }
}

income.onChangeDateStart = function(val){
    if (val.getTime()>income.DateEnd().getTime()){
        income.DateEnd(val)
    }
}

income.init = function () {
    // income.setDate()
    income2.renderFilterBy()
    income.getData(function () {
        income.rendergridSales()
        income.rendergridIncomeRetur()
        // income.rendergridIncomeStatement1_5()
        income.rendergridOperatingExpenses()
        income.rendergridOtherIncome()
        income.rendergridOtherExpense()
        income.removeRowGroup()
        income.getDataTax(function () {
            income.renderGridTAX()
            income.getDateNow()
        })
    })

}

$(function () {
    income.init()
})