var income2={}
income2.filterByValue= ko.observable("Detail")
income2.buttonExportExcel = ko.observable(false)
income2.renderFilterBy= function(){
    var Data = [
    {"value":"Detail", "text":"Detail"},
    {"value":"Periode", "text":"Period"}
    ]
    $("#filterBy").kendoDropDownList({
        filter: "contains",
        dataTextField: "text",
        dataValueField: "value",
        dataSource: Data,
        // index: 1,
        // optionLabel:'Select one',
        change:function(e){
            var dataitem = this.dataItem();
            income2.filterByValue(dataitem.value)

        },
    });
}
income2.getDataPeriode = function (callback) {
    model.Processing(true)
    var url = "/financial/getdatatrialincomeperiode"
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var param = {}
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
    // console.log(param.DateStart)
    // console.log(param.DateEnd)
    // console.log(param.TextSearch)
    ajaxPost(url, param, function (res) {
        if (res.IsError === "true") {
            swal("Error!", res.Message, "error")
            return
        }
        // console.log("Data ### =>",res.DataItem);
        var monthyear = []
        for (var d = new Date(dateStart); d <= new Date(dateEnd); d.setMonth(d.getMonth() + 1)) {
            // console.log("6easdasfafhsdb")
            monthyear.push({
               M : parseInt(moment(d).format("M")),
               Y : parseInt(moment(d).format("YYYY")),
               my :moment(d).format("MMMYYYY")
           })
        }
        // console.log("Date : ",monthyear)
        // var data = []
        var DataAmount =[]
        _.forEach(res.Data.DataAmount , function(v,i){
            var m = v.Month
            var y = v.Year 
            var dateMy = new Date(y,m-1, 01)
            var my = moment(dateMy).format("MMMYYYY")
            DataAmount.push({
                Acc_Code: v.Acc_Code,
                Acc_Name :v.Acc_Name,
                MonthYear : my,
                Amount : v.Amount
            })
        })
        var dataAmont2 =_.chain(DataAmount).groupBy("Acc_Code").toPairs().map(function (currentItem) {
            return _.zipObject(["Acc_Code","DataItem"], currentItem);
        }).value();
        // console.log("Data Amount:",DataAmount, "Data Amont2:",dataAmont2)
        var dataAcc = []
        _.forEach(res.Data.DataAcc, function(v,i){
            
            var name = v.Acc_Code
            var each = _.find(dataAmont2, function(vv){return vv.Acc_Code==name})
            if (each!= undefined) {
                var itemMy = {}
                var totalAmount = 0
                _.forEach(monthyear, function(key){
                    var eachItem = _.find(each.DataItem, function(vv){return vv.MonthYear==key.my}) 
                    // console.log("data =>", eachItem.DataItem)
                    if(eachItem==undefined){
                        itemMy[key.my] = 0 
                    }else{
                        itemMy[key.my] =eachItem.Amount
                        totalAmount = totalAmount+ eachItem.Amount
                    }
                })
                dataAcc.push({
                    Acc_Code:v.Acc_Code,
                    Main_Acc_Code : v.Main_Code,
                    Acc_Name:v.Acc_Name,
                    DataItem : itemMy,
                    TotalAmount: totalAmount
                })
            }else{
                var itemMy = {}
                _.forEach(monthyear, function(key){
                    itemMy[key.my] = 0 
                })
                // console.log("Itemkey ->", itemMy)
                dataAcc.push({
                    Acc_Code:v.Acc_Code,
                    Main_Acc_Code : v.Main_Code,
                    Acc_Name:v.Acc_Name,
                    DataItem : itemMy,
                    TotalAmount: 0
                })
                // console.log("DataItem =>", itemMy)
            }
            // console.log("DataAcc =>", dataAcc)
        })
        income.dataMasterIncome(dataAcc)
        model.Processing(false)
        callback()
    }, function () {
        swal("Error!", "Unknown error, please try again", "error")
    })
}
income2.rendergridSalesPeriod = function () {
    var Data5000 = _.filter(income.dataMasterIncome(), function(x){
        return x.Acc_Code== 5110
    });
    for (i in Data5000) {
        Data5000[i].Parent = "PENJUALAN"
    }
    var data = Data5000
    // console.log("salees",data)
    var columns = [{
        field: 'Acc_Code',
        title: 'Acc Code',
        width: 50,
        hidden: true,
    },{
        field: 'Acc_Name',
        title: 'Account Name',
        width: 150,
        // footerTemplate: "<div style='text-align:left; font-size: 15px;'>TOTAL SALES AND REVENUE:</div>"
    },{
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #",
        hidden: true,
    }]
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var Aggregate = []
    for (var d = new Date(dateStart); d <= new Date(dateEnd); d.setMonth(d.getMonth() + 1)) {
        var namefield = "DataItem." +moment(d).format("MMMYYYY")
        var date1 = moment(d).format("MMMYYYY")
        // console.log("=>", namefield)
        columns.push({
            title: moment(d).format("MMM YYYY"),
            field: namefield,
            template: "#=ChangeToRupiah("+namefield+")#",
            width: 100,
            attributes: {
                style: "text-align:right;"
            },
            // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
            footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= income2.calculateGrossEndingPeriod('"+namefield+"') #</div>",
        })
        Aggregate.push({
            field: namefield,
            aggregate: "sum"
        })
    }
    columns.push({
        title: "TOTAL",
        field: "TotalAmount",
        width: 100,
        attributes: {
            style: "text-align:right;"
        },
        hidden: true,
        template: "#=ChangeToRupiah(TotalAmount)#",
        footerTemplate: "<div class='align-right'>#=ChangeToRupiah(sum)#</div>",
    })
    Aggregate.push({
        field: "TotalAmount",
        aggregate: "sum"
    })

    $('#gridIncomeSales').html("")
    $('#gridIncomeSales').kendoGrid({
        dataSource: {
            data: data,
            aggregate: Aggregate,
            group: {
                field: "Parent",
                dir: "asc"
            }
        },
        scrollable: true,
        columns: columns
    })
    $("#gridIncomeSales .k-grid-footer").css('display', 'none');
}
income2.calculateGrossEndingPeriod = function (namefield) {
    var fieldS = namefield.split(".")
    var data = $("#gridIncomeSales").data('kendoGrid').dataSource.options.data
    console.log("Calculate =>", data)
    var DataSum = 0.0
    var DataMin = 0.0
    // console.log(data)
    for (i in data) {
        var each = data[i]
            // console.log(each.DataItem[fieldS[1]])
        if (data[i].Acc_Code == 5600) {
            DataMin += each.DataItem[fieldS[1]]
        } else {
            DataSum += each.DataItem[fieldS[1]]
        }
    }
    var Results = DataSum - DataMin
    // income.koGrossEndingValue(Results)
    return kendo.toString(Results, 'n')
    // return ""
}
income2.rendergridIncomeRetur = function(){
    var data = _.filter(income.dataMasterIncome(), function(x){
        return x.Acc_Code == 5210 || x.Acc_Code== 5211 || x.Acc_Code==5212  
    });
    // var totalIncomeRetur = _.sumBy(data, function(o) { return o.Ending; })
    for (i in data) {
        data[i].Parent = "HARGA POKOK"
    }
    var columns = [{
        field: 'Acc_Code',
        title: 'Acc Code',
        width: 50,
        hidden: true,
    }, {
        field: 'Acc_Name',
        title: 'Account Name',
        width: 150,
        footerTemplate: "<div style='text-align:left; font-size: 15px;'>TOTAL HARGA POKOK:</div><div style='text-align:left; font-size: 15px;'>LABA KOTOR:</div>"
    }, {
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #",
        hidden: true,
    }]
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var Aggregate = []
    for (var d = new Date(dateStart); d <= new Date(dateEnd); d.setMonth(d.getMonth() + 1)) {
        var namefield = "DataItem." +moment(d).format("MMMYYYY")
        var date1 = moment(d).format("MMMYYYY")
        console.log("hpp ", namefield)
        columns.push({
            title: moment(d).format("MMM YYYY"),
            field: namefield,
            template: "#=ChangeToRupiah("+namefield+")#",
            width: 100,
            attributes: {
                style: "text-align:right;"
            },
            // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",#= ChangeToRupiah(income2.sumSalesEndingLabaKotorPeriod('"+namefield+"', sum)) #, #= ChangeToRupiah(income2.sumTotalHPP("+namefield+")) #
            footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income2.sumSalesEndingLabaKotorPeriod('"+namefield+"', sum)) #</div><div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income2.sumTotalHPP('"+namefield+"')) #</div>",
        })
        Aggregate.push({
            field: namefield,
            aggregate: "sum"
        })
    }
    columns.push({
        title: "TOTAL",
        field: "TotalAmount",
        width: 100,
        attributes: {
            style: "text-align:right;"
        },
        hidden: true,
        template: "#=ChangeToRupiah(TotalAmount)#",
        // footerTemplate: "<div class='align-right'>#=ChangeToRupiah(sum)#</div>",
    })
    Aggregate.push({
        field: "TotalAmount",
        aggregate: "sum"
    })
    console.log("tes data hpp", data)

    $('#gridIncomeRetur').html("")
    $('#gridIncomeRetur').kendoGrid({
        dataSource: {
            data: data,
            aggregate: Aggregate,
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
income2.sumPenjualanBersih= function(namefield, sum){
    var totalPenjualan = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates()
    var result =  totalPenjualan[namefield].sum - sum
    return result
}
income2.rendergridIncomeStatement1_5Period = function () {
    var Data5000 = _.filter(income.dataMasterIncome(), function(x){
        return x.Acc_Code== 5210 || x.Main_Acc_Code==5200  
    });
    for (i in Data5000) {
        Data5000[i].Parent = "Income Statement"
    }
    var data = Data5000
    // console.log("salees",data)
    var columns = [{
        field: 'Acc_Code',
        title: 'Acc Code',
        width: 50,
        hidden: true,
    },{
        field: 'Acc_Name',
        title: 'Account Name',
        width: 150,
        footerTemplate: "<div style='text-align:left; font-size: 15px;'>TOTAL HARGA POKOK PENJUALAN:</div><div style='text-align:left; font-size: 15px;'>TOTAL LABA KOTOR:</div>"
    },{
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #",
        hidden: true,
    }]
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var Aggregate = []
    for (var d = new Date(dateStart); d <= new Date(dateEnd); d.setMonth(d.getMonth() + 1)) {
        var namefield = "DataItem." +moment(d).format("MMMYYYY")
        var date1 = moment(d).format("MMMYYYY")
        // console.log(namefield)
        columns.push({
            title: moment(d).format("MMM YYYY"),
            field: namefield,
            template: "#=ChangeToRupiah("+namefield+")#",
            width: 100,
            attributes: {
                style: "text-align:right;"
            },
            // #= ChangeToRupiah(income2.sumSalesEndingLabaKotorPeriod('"+namefield+"')) #
            // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
            footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income2.sumTotalHPP('"+namefield+"')) #</div><div style='text-align:right; font-size: 15px;'></div>",
        })
        Aggregate.push({
            field: namefield,
            aggregate: "sum"
        })
    }
    columns.push({
        title: "TOTAL",
        field: "TotalAmount",
        width: 100,
        attributes: {
            style: "text-align:right;"
        },
        hidden: true,
        template: "#=ChangeToRupiah(TotalAmount)#",
        footerTemplate: "<div class='align-right'>#=ChangeToRupiah(sum)#</div>",
    })
    Aggregate.push({
        field: "TotalAmount",
        aggregate: "sum"
    })

    $('#gridIncomeStatement1_5').html("")
    $('#gridIncomeStatement1_5').kendoGrid({
        dataSource: {
            data: data,
            aggregate: Aggregate,
            group: {
                field: "Parent",
                dir: "asc"
            }
        },
        scrollable: true,
        columns: columns
    })
    $("#gridIncomeStatement1_5 > div.k-grid-header > div > table > thead > tr").hide()
}
income2.sumTotalHPP = function(fieldName){
    var TotalHPP = $("#gridIncomeRetur").data("kendoGrid").dataSource.aggregates()
    var PotonganDanRetur = $("#gridIncomeRetur").data("kendoGrid").dataSource.aggregates()
    // var totalPotonganDanRetur = income2.sumPenjualanBersih(fieldName, PotonganDanRetur[fieldName].sum)
    
    return TotalHPP[fieldName].sum
}
income2.sumSalesEndingLabaKotorPeriod = function(fieldName){
    var PotonganDanRetur = $("#gridIncomeRetur").data("kendoGrid").dataSource.aggregates()
    var TotalHPP=  income2.sumTotalHPP(fieldName)
    var PenjualanBersih = income2.sumPenjualanBersih(fieldName, PotonganDanRetur[fieldName].sum)
    // console.log(PotonganDanRetur[fieldName].sum,TotalHPP)
    return PenjualanBersih-TotalHPP
    // var v5210 = $("#gridIncomeStatement1_5").data("kendoGrid").dataSource.aggregates()
    // var totalEnding = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates()
    // var result = totalEnding[fieldName].sum- v5210[fieldName].sum
    // return ChangeToRupiah(result)
}


income2.rendergridOperatingExpensesPeriod = function() {
    var Data6001 = _.filter(income.dataMasterIncome(), function(x) {
        return x.ACC_Code!=6160&&x.ACC_Code!=6161&&x.ACC_Code!=6999&& x.Main_Acc_Code == 6000
    });
    console.log("Data 6001 : ",Data6001)
    for (i in Data6001) {
        Data6001[i].Parent = "BIAYA UMUM DAN ADMINISTRASI"
    }
    var dataPeriod = Data6001
    var columns = [{
        field: 'Acc_Code',
        title: 'Acc Code',
        width: 50,
        hidden: true
    },{
        field: 'Acc_Name',
        title: 'Account name',
        width: 150,
        footerTemplate: "<div style='text-align:left; font-size: 15px;'>TOTAL BIAYA UMUM DAN ADMINISTRASI:</div><div style='text-align:left; font-size: 15px;'>LABA USAHA:</div>"
    },{
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #",
        hidden: true,
    }]
    var dateStart = $('#dateStart').data('kendoDatePicker').value()
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value()
    var Aggregate = []
    for (var d = new Date(dateStart); d <= new Date(dateEnd); d.setMonth(d.getMonth() + 1)) {
        var namefield = "DataItem."+moment(d).format("MMMYYYY")
        var date2 = moment(d).format("MMMYYYY")
        columns.push({
            title: moment(d).format("MMM YYYY"),
            field: namefield,
            template: "#=ChangeToRupiah("+namefield+")#",
            width: 100,
            attributes: {
                style: "text-align:right;"
            },
            footerTemplate: "<div style='text-align=right;font-size=15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align=right;font-size: 15px'>#= ChangeToRupiah(income2.netProfitPeriod('"+namefield+"')) #</div>",
        })
        Aggregate.push({
            field: namefield,
            aggregate: "sum"
        })
    }
    columns.push({
        title: "TOTAL",
        field: "TotalAmount",
        width: 100,
        attributes: {
            style: "text-align:right;"
        },
        hidden: true,
        template: "#=ChangeToRupiah(TotalAmount)#",
        footerTemplate: "<div class='align-right'>#=ChangeToRupiah(sum)#</div>",
    })
    Aggregate.push({
        field: "TotalAmount",
        aggregate: "sum"
    })

    $("#gridIncomeOperatingExpenses").kendoGrid({
        dataSource: {
            data: dataPeriod,
            aggregate: Aggregate,
            sort: {
                field: "Acc_Code",
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

income2.netProfitPeriod = function (fieldName) {
    var labakotor = income2.sumSalesEndingLabaKotorPeriod(fieldName)
    var IncomeOperatingExpenses = $("#gridIncomeOperatingExpenses").data("kendoGrid").dataSource.aggregates()
    return labakotor - IncomeOperatingExpenses[fieldName].sum

    // var v5210 = $("#gridIncomeStatement1_5").data("kendoGrid").dataSource.aggregates()
    // var totalEnding = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates()
    // var LabaKotor = totalEnding[fieldName].sum- v5210[fieldName].sum

    // var dataItemOE = $("#gridIncomeOperatingExpenses").data("kendoGrid").dataSource.aggregates()
    // var totalEndingOE = dataItemOE[fieldName].sum
    // // var dataGrossEnding = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates()
    // // var totalGrossEnding = dataGrossEnding[fieldName].sum
    // var Results = LabaKotor- totalEndingOE
    // // console.log(totalEndingOE,totalGrossEnding)
    // return ChangeToRupiah(Results)
}

income2.rendergridOtherIncomePeriod = function () {
    var Data7001 = _.filter(income.dataMasterIncome(), function(x) {
        return x.Acc_Code == 7100 || x.Acc_Code == 7200 || x.Acc_Code == 7999
    })
    console.log("Data 7001 =>", Data7001)
    var dataPeriod = Data7001
    for(i in Data7001) {
        Data7001[i].Parent = "PENDAPATAN/(BIAYA) DI LUAR USAHA"
    }
    var columns = [{
        field: 'Acc_Code',
        title: 'Acc Code',
        width: 50,
        hidden: true,
    },{
        field: 'Acc_Name',
        title: 'Account Name',
        width: 150,
        footerTemplate: "<div style='text-align:left; font-size: 15px;'>TOTAL PENDAPATAN DI LUAR USAHA:</div>"
    },{
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #",
        hidden: true,
    }]
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var Aggregate = []
    for (var d = new Date(dateStart); d <= new Date(dateEnd); d.setMonth(d.getMonth() + 1)) {
        var namefield = "DataItem." +moment(d).format("MMMYYYY")
        var date1 = moment(d).format("MMMYYYY")
        // console.log(namefield)
        columns.push({
            title: moment(d).format("MMM YYYY"),
            field: namefield,
            template: "#=ChangeToRupiah("+namefield+")#",
            width: 100,
            attributes: {
                style: "text-align:right;"
            },
            footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        })
        Aggregate.push({
            field: namefield,
            aggregate: "sum"
        })
    }
    columns.push({
        title: "TOTAL",
        field: "TotalAmount",
        width: 100,
        attributes: {
            style: "text-align:right;"
        },
        hidden: true,
        template: "#=ChangeToRupiah(TotalAmount)#",
        footerTemplate: "<div class='align-right'>#=ChangeToRupiah(sum)#</div>",
    })
    Aggregate.push({
        field: "TotalAmount",
        aggregate: "sum"
    })
    $('#gridIncomeOtherIncome').kendoGrid({
        dataSource: {
            data: dataPeriod,
            aggregate: Aggregate,
            sort: {
                field: "Acc_Code",
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
    $("#gridIncomeOtherIncome > div.k-grid-header > div > table > thead > tr").hide()
    }

income2.rendergridOtherExpensePeriod = function () {
    var Data8001 = _.filter(income.dataMasterIncome(), function(x) {
        return  x.Acc_Code == 8100 || x.Acc_code == 8200 || x.Acc_Code == 8300 || x.Acc_Code == 8400 || x.Acc_Code == 8500
    })
    console.log("Data 8001 =>", Data8001)
    for (i in Data8001) {
        Data8001[i].Parent = "BIAYA DI LUAR USAHA"
    }
    var dataPeriod = Data8001
    var columns = [{
        field: 'Acc_Code',
        title: 'Acc Code',
        width: 50,
        hidden: true,
    },{
        field: 'Acc_Name',
        title: 'Account Name',
        width: 150,
        footerTemplate: "<div style='text-align:left;font-size:15px;'>TOTAL BIAYA DI LUAR USAHA:</div><div style='text-align:left;font-size:15px;'>TOTAL BIAYA DAN PENDAPATAN DI LUAR USAHA:</div><div style='text-align:left;font-size:15px;'>LABA BERSIH SEBELUM PAJAK:</div>"
    },{
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #",
        hidden: true
    }]
    var dateStart = $('#dateStart').data('kendoDatePicker').value()
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value()
    var Aggregate = []
    for (var d = new Date(dateStart); d <= new Date(dateEnd); d.setMonth(d.getMonth() + 1)) {
        var namefield = "DataItem."+moment(d).format("MMMYYYY")
        var date1 = moment(d).format("MMMYYYY")
        columns.push({
            title: moment(d).format("MMM YYYY"),
            field: namefield,
            template: "#=ChangeToRupiah("+namefield+")#",
            width: 100,
            attributes: {
                style: "text-align:right;"
            },
            footerTemplate : "<div style='text-align:right;font-size:15px;'>#= ChangeToRupiah(sum) #</div><div style='text-align:right;font-size:15px;'>#=income2.calculateRevenueExpensesPeriod('"+namefield+"')#</div><div style='text-align:right;font-size:15px;'>#= ChangeToRupiah(income2.calculateEarningBeforTaxPeriod('"+namefield+"')) #</div>",
        })
        Aggregate.push({
            field: namefield,
            aggregate: "sum"
        })
    }
    columns.push({
        title: "TOTAL",
        field: "TotalAmount",
        width: 100,
        attributes: {
            style: "text-align:right;"
        },
        hidden: true,
        template: "#=ChangeToRupiah(TotalAmount)#",
        footerTemplate: "<div class='align-right'>#=ChangeToRupiah(sum)#</div>",
    })
    Aggregate.push({
        field: "TotalAmount",
        aggregate: "sum"
    })
    $("#gridIncomeOtherExpenses").kendoGrid({
        dataSource: {
            data: dataPeriod,
            aggregate: Aggregate,
            sort: {
                field: 'Acc_Code',
                dir: 'asc'
            },
            group: {
                field: 'Parent',
                dir: 'asc'
            }
        },
        scrollable: true,
        columns: columns
    })

    $("#gridIncomeOtherExpenses > div.k-grid-header > div > table > thead > tr").hide()
}

income2.calculateRevenueExpensesPeriod = function(fieldName){
    var revenue = $("#gridIncomeOtherIncome").data("kendoGrid").dataSource.aggregates()
    var totalRevenue= revenue[fieldName].sum
    var expenses = $("#gridIncomeOtherExpenses").data("kendoGrid").dataSource.aggregates()
    var totalExpenses = expenses[fieldName].sum
    var Results = totalRevenue - totalExpenses
    return ChangeToRupiah(Results)
}
income2.calculateEarningBeforTaxPeriod = function(fieldName){
    // // var dataGrossEnding = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates()
    // // var totalGrossEnding = dataGrossEnding[fieldName].sum
    // var v5210 = $("#gridIncomeStatement1_5").data("kendoGrid").dataSource.aggregates()
    // var totalEnding = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates()
    // var LabaKotor = totalEnding[fieldName].sum- v5210[fieldName].sum
    // var dataItemOE = $("#gridIncomeOperatingExpenses").data("kendoGrid").dataSource.aggregates()
    // var totalEndingOE = dataItemOE[fieldName].sum
    // var ResultsNetProfit = LabaKotor -totalEndingOE
    // // console.log(totalGrossEnding,totalExpenses)
    // //
    // var Expenses = $("#gridIncomeOtherExpenses").data("kendoGrid").dataSource.aggregates()
    // var totalExpenses = Expenses[fieldName].sum
    // var Revenue = $("#gridIncomeOtherIncome").data("kendoGrid").dataSource.aggregates()
    // var totalRevenue= Revenue[fieldName].sum
    // var ResultsRevExp = totalRevenue - totalExpenses
    // var Results = ResultsNetProfit + ResultsRevExp
    // return ChangeToRupiah(Results)
    var netProfitPeriod = income2.netProfitPeriod(fieldName)
    var Expenses = $("#gridIncomeOtherExpenses").data("kendoGrid").dataSource.aggregates()
    var totalExpenses = Expenses[fieldName].sum
    var Revenue = $("#gridIncomeOtherIncome").data("kendoGrid").dataSource.aggregates()
    var totalRevenue= Revenue[fieldName].sum
    var ResultsRevExp = totalRevenue - totalExpenses
    var Results = netProfitPeriod + ResultsRevExp
    return Results
}
income2.getDataTaxPeriod = function (callback) {
    model.Processing(true)
    var url = "/financial/getdatataxperiod"
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var param = {}
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
        if (res.IsError === "true") {
            swal("Error!", res.Message, "error")
            return
        }
        var monthyear = []
        for (var d = new Date(dateStart); d <= new Date(dateEnd); d.setMonth(d.getMonth() + 1)) {
            // console.log("6easdasfafhsdb")
            monthyear.push({
               M : parseInt(moment(d).format("M")),
               Y : parseInt(moment(d).format("YYYY")),
               my :moment(d).format("MMMYYYY")
           })
        }
        // var data = []
        var DataAmount =[]
        _.forEach(res.Data.DataAmount , function(v,i){
            var m = v.Month
            var y = v.Year 
            var dateMy = new Date(y,m-1, 01)
            var my = moment(dateMy).format("MMMYYYY")
            DataAmount.push({
                Acc_Code: v.Acc_Code,
                Acc_Name :v.Acc_Name,
                MonthYear : my,
                Amount : v.Amount
            })
        })
        var dataAmont2 =_.chain(DataAmount).groupBy("Acc_Code").toPairs().map(function (currentItem) {
            return _.zipObject(["Acc_Code","DataItem"], currentItem);
        }).value();
        // console.log(DataAmount, dataAmont2)
        var dataAcc = []
        _.forEach(res.Data.DataAcc, function(v,i){
            var name = v.Acc_Code
            var each = _.find(dataAmont2, function(vv){return vv.Acc_Code==name})
            if (each!= undefined) {
                var itemMy = {}
                var totalAmount = 0
                _.forEach(monthyear, function(key){
                    var eachItem = _.find(each.DataItem, function(vv){return vv.MonthYear==key.my}) 
                    if(eachItem==undefined){
                        itemMy[key.my] = 0 
                    }else{
                        itemMy[key.my] =eachItem.Amount
                        totalAmount = totalAmount+ eachItem.Amount
                    }
                })
                dataAcc.push({
                    Acc_Code:v.Acc_Code,
                    Main_Acc_Code : v.Main_Code,
                    Acc_Name:v.Acc_Name,
                    DataItem : itemMy,
                    TotalAmount: totalAmount
                })
            }else{
                var itemMy = {}
                _.forEach(monthyear, function(key){
                    itemMy[key.my] = 0 
                })
                dataAcc.push({
                    Acc_Code:v.Acc_Code,
                    Main_Acc_Code : v.Main_Code,
                    Acc_Name:v.Acc_Name,
                    DataItem : itemMy,
                    TotalAmount: 0
                })
            }
        })
        // console.log(dataAcc)
        income.dataTax(dataAcc)
        model.Processing(false)
        callback()
    }, function () {
        swal("Error!", "Unknown error, please try again", "error")
    })
}
income2.renderGridTAXPeriod = function () {
    var dataPeriod = income.dataTax()
    for (i in dataPeriod) {
        dataPeriod[i].Parent = ""
    }
    var columns = [{
        field: 'Acc_Name',
        title: 'Account Name',
        width: 150,
        footerTemplate: "<div style='text-align:left; font-size: 15px;'>LABA BERSIH SETELAH PAJAK:</div>"
    }, {
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #",
        hidden: true,
    }]
    var dateStart = $('#dateStart').data('kendoDatePicker').value()
    var dateEnd = $("#dateEnd").data('kendoDatePicker').value()
    var Aggregate = []
    for(var d = new Date(dateStart); d <= new Date(dateEnd); d.setMonth(d.getMonth() + 1)) {
        var namefield = "DataItem."+moment(d).format("MMMYYYY")
        var date1 = moment(d).format("MMMYYYY")
        columns.push({
            title: moment(d).format("MMM YYYY"),
            field: namefield,
            template: "#=ChangeToRupiah("+namefield+")#",
            width: 100,
            attributes: {
                style: "text-align:right;"
            },
            footerTemplate: "<div style='text-align:right;font-size:15px;'>#=ChangeToRupiah(income2.calculateEarningAfterTax2('"+namefield+"'))#</div>",
        })
        Aggregate.push({
            field: namefield,
            aggregate: "sum"
        })
    }
    columns.push({
        title: "TOTAL",
        field: 'TotalAmount',
        width: 100,
        attributes: {
            style: "text-align:right;"
        },
        hidden: true,
        template: "#=ChangeToRupiah(TotalAmount)#",
        footerTemplate: "<div class='align-right'>#=ChangeToRupiah(sum)#</div>"
    })
    Aggregate.push({
        field: "TotalAmount",
        aggregate: "sum"
    })
    // $('#gridTAX').html('')
    $('#gridTAX').kendoGrid({
        dataSource: {
            data: dataPeriod,
            aggregate: Aggregate,
            sort: {
                field: 'Acc_Code',
                dir: 'asc'
            },
            group: {
                field: 'Parent',
                dir: 'asc'
            }
        },
        scrollable: true,
        columns: columns
    })
    $("#gridTAX > div.k-grid-header > div > table > thead > tr").hide()
}
income2.calculateEarningAfterTax2 = function (fieldName) {
    // console.log(value)
    // var dataGrossEnding = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates()
    // var totalGrossEnding = dataGrossEnding[fieldName].sum
    // var v5210 = $("#gridIncomeStatement1_5").data("kendoGrid").dataSource.aggregates()
    // var totalEnding = $("#gridIncomeSales").data("kendoGrid").dataSource.aggregates()
    // var LabaKotor = totalEnding[fieldName].sum- v5210[fieldName].sum
    // var dataItemOE = $("#gridIncomeOperatingExpenses").data("kendoGrid").dataSource.aggregates()
    // var totalEndingOE = dataItemOE[fieldName].sum
    // var ResultsNetProfit = LabaKotor -totalEndingOE
    // // console.log(totalGrossEnding,totalExpenses)
    // //
    // var Expenses = $("#gridIncomeOtherExpenses").data("kendoGrid").dataSource.aggregates()
    // var totalExpenses = Expenses[fieldName].sum
    // var Revenue = $("#gridIncomeOtherIncome").data("kendoGrid").dataSource.aggregates()
    // var totalRevenue= Revenue[fieldName].sum
    // var ResultsRevExp = totalRevenue - totalExpenses
    // var ResultsBefore = ResultsNetProfit + ResultsRevExp
    // // tax
    // var dataTax = $("#gridTAX").data("kendoGrid").dataSource.aggregates()
    // var totalTax = dataTax[fieldName].sum
    // // console.log(ResultsBefore,totalTax)
    // var Results = ResultsBefore - totalTax
    // return ChangeToRupiah(Results)
    var totalBfrTax = income2.calculateEarningBeforTaxPeriod(fieldName)
    var tax = $("#gridTAX").data("kendoGrid").dataSource.aggregates()
    return totalBfrTax -tax[fieldName].sum
}

income2.exportExcelPeriod = function(){
    model.Processing(true)
    var url = "/financial/exporttoexcelperiod"
    var dateStart = new Date(kendo.toString($('#dateStart').data('kendoDatePicker').value(),"yyyy-MM-dd"))
    var dateEnd = new Date(kendo.toString($('#dateEnd').data('kendoDatePicker').value(),"yyyy-MM-dd"))
    var param = {
        Data :_.sortBy(income.dataMasterIncome(), [function(o) { return o.Acc_Code; }]),
        DataTax :_.sortBy(income.dataTax(), [function(o) { return o.Acc_Code; }]),
        DateStart: dateStart,
        DateEnd: dateEnd,
    }
    ajaxPost(url, param, function (e) {
        var pom = document.createElement('a');
        pom.setAttribute('href', "/res/docs/report/excel/" + e.Data);
        pom.setAttribute('download', e.Data);
        pom.click();
        model.Processing(false)
    }, function () {
        swal("Error!", "Unknown error, please try again", "error")
    })
}