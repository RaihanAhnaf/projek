var trial2 = {}
trial2.dataAktiva = ko.observable()
trial2.dataPasiva = ko.observable()
trial2.dataIncomePeriode = ko.observable()
trial2.dataTax = ko.observable()
trial2.renderGridActivaPeriod = function() {
    var Data = trial2.dataAktiva()

    // for(i in Data) {
    //     Data[i].Parent = "Aktiva"
    var columns = [
        // {
        //     field: "Parent",
        //     title: "Parent",
        //     groupHeaderTemplate: "#= value #",
        //     hidden: true
        // },
        {
        field: "Acc_Code",
        title: "Acc Code",
        width: 50
    },
     {
        field: "Acc_Name",
        title: "Account Name",
        width: 70,
        footerTemplate: "<div style='text-align:center; font-size: 15px;'>Total Aktiva:</div>"
    }
]
    var dateStart = $('#dateStart').data('kendoDatePicker').value()
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value()
    var Aggregate = []
    for(var d = new Date(dateStart); d <= new Date(dateEnd);d.setMonth(d.getMonth() + 1)) {
        var namefield = "DataItem."+moment(d).format("MMMYYYY")+"."
        columns.push(
         {
            field: namefield+"begining",
            title: moment(d).format("MMMYYYY")+'-Begining',
            // template: "#=ChangeToRupiah(sum)#",
            footerTemplate: "<div style='text-align:right;'>#=ChangeToRupiah(sum)#</div>",
            attributes: {
                style: "text-align:right;"
            },
            width: 100
        }, {
            field: namefield+"debet",
            title: moment(d).format("MMMYYYY")+'-Debet',
            // template: "#=ChangeToRupiah(sum)#",
            footerTemplate: "<div style='text-align:right;font-size:15px;'>#=ChangeToRupiah(sum)#</div>",
            attributes: {
                style: "text-align:right;"
            },
            width: 100
        }, {
            field: namefield+"credit",
            title: moment(d).format("MMMYYYY")+'-Credit',
            // template: "#=ChangeToRupiah(sum)#",
            width: 100,
            footerTemplate: "<div style='text-align:right;font-size:15px;'>#=ChangeToRupiah(sum)#</div>",
            attributes: {
                style: "text-align:right;"
            }
        }, {
            field: namefield+"ending",
            title: moment(d).format("MMMYYYY")+'-Ending',
            // template: "#=ChangeToRupiah(sum)#",
            width: 100,
            footerTemplate: "<div style='text-align:right;font-size:15px'>#=ChangeToRupiah(sum)#</div>",
            attributes: {
                style: "text-align:right;"
            }
        })
        Aggregate.push({
            field: namefield+"begining",
            aggregate: "sum"
        }, {
            field: namefield+"debet",
            aggregate: "sum"
        }, {
            field: namefield+"credit",
            aggregate: "sum"
        }, {
            field: namefield+"ending",
            aggregate: "sum"
        })
    }
    // console.log(Aggregate)
    columns.push({
        title: "TOTAL",
        field: "TotalAmount",
        width: 100,
        attributes: {
            style: "text-align:right;"
        },
        hidden: true,
        // template: "#=ChangeToRupiah(TotalAmount)#"
        // footerTemplate: "<div class='align-right'>#=ChangeToRupiah(sum)#</div>",
    })

    $("#gridActiva").html("");
    $("#gridActiva").kendoGrid({
        dataSource: {
            data: Data,
            aggregate: Aggregate
            // group: [{
            //     dir: "asc",
            //     field: "Acc_Code"
            // }]
        },
        height: 400,
        scrollable: true,
        columns: columns,
        columnmenu: true
    });

    // $("#gridActiva > div.k-grid-header > div > table > thead > tr").hide()
}

trial2.getDataActiva = function (callback) {
    model.Processing(true)
    var url = "/financial/getdataactivaperiode"
    var dateStart = $('#dateStart').data('KendoDatePicker').value()
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value()
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

trial2.getDataActivaPeriod = function (callback){
    model.Processing(true)
    var url = "/financial/getdataactivaperiode"
    var dateStart = $("#dateStart").data("kendoDatePicker").value();
    // var dateEnd = $("#dateEnd").data("KendoDatePicker").value()
    var dateEnd = $("#dateEnd").data("kendoDatePicker").value();
    // console.log("DateStart :", dateStart, " - DateEnd :", dateEnd)
    var param = {}
    if (trial.koFilterIsActive() == true) {
        param = {
            DateStart: dateStart,
            DateEnd: dateEnd,
            TextSearch: trial.textSearch(),
            Filter: true
        }
    }else {
        param = {
            Filter: false,
        }
    }

    ajaxPost(url, param, function (res) {
        if (res.IsError === "true") {
            swal("Error!", res.Message, "error")
            return
        }
        var dataField = []
        _.forEach(res.Data.DataAcc, function (v) {
            dataField.push({
                beginingVal: v.Amount_Credit,
                creditVal: v.Amount_Credit,
                debetVal: v.Amount_Debet,
                endingVal: v.Amount_Ending
            })
        })
        var monthYear = []
        for (var d = new Date(dateStart);d <= new Date(dateEnd);d.setMonth(d.getMonth() + 1)) {
            var keyDate = moment(d).format("MMMYYYY")
            if (dataField.beginingVal == undefined) {
                dataField.beginingVal = 0
            }
            if (dataField.creditVal == undefined) {
                dataField.creditVal = 0
            }
            if (dataField.debetVal == undefined) {
                dataField.debetVal = 0
            }
            if (dataField.endingVal == undefined) {
                dataField.endingVal = 0
            }
            monthYear.push({
                M : parseInt(moment(d).format("M")),
                Y : parseInt(moment(d).format("YYYY")),
                MY : keyDate,
                my :  {
                         begining: dataField.beginingVal,
                         credit:dataField.creditVal,
                         debet: dataField.debetVal,
                         ending: dataField.endingVal
                   }
            })
        }
        var DataAmount = []
        _.forEach(res.Data.DataAmount, function(v,i) {
            var m = v.Month
            var y = v.Year
            var dateMy = new dateEnd(y,m-1, 01)
            var my = moment(dateMy).format("MMMYYYY")
            DataAmount.push({
                Acc_Code: v.Acc_Code,
                Acc_Name: v.Acc_Name,
                MonthYear: my,
                Amount: v.Amount
            })
        })
        var dataAmount2 = _.chain(DataAmount).groupBy("Acc_Code").toPairs().map(function (curentItem) {
            return _.zipObject(["Acc_Code", "DataItem"], curentItem);
        }).value();
        var dataAcc = []
        _.forEach(res.Data.DataAcc, function(v,i) {
            var name = v.Acc_Code
            var each = _.find(dataAmount2.DataItem, function(vv){return vv.MonthYear==name})
            // console.log("Each :", each)
            if(each != undefined) {
                var itemMy = {}
                var totalAmount = 0
                _.forEach(monthYear, function(key) {
                    var eachItem = _.find(each.DataItem, function(vv){return vv.MonthYear==key.my})
                    if (eachItem == undefined) {
                        itemMy[key.MY] = 0
                    }else {
                        itemMy[key.my] = eachItem.Amount
                        totalAmount = totalAmount + eachItem
                    }
                })
                dataAcc.push({
                    Acc_Code: v.Acc_Code,
                    Main_Acc_Code: v.Main_Code,
                    Acc_Name: v.Acc_name,
                    DataItem: itemMy,
                    TotalAmount: totalAmount
                })
            }else {
                var itemMy = {}
                _.forEach(monthYear, function(key) {
                    itemMy[key.MY] = key.my
                })
                dataAcc.push({
                    Acc_Code: v.Acc_Code,
                    Main_Acc_Code: v.Main_Code,
                    Acc_Name: v.Acc_Name,
                    DataItem: itemMy,
                    TotalAmount: 0
                })
            }
        })
        trial2.dataAktiva(dataAcc)
        model.Processing(false)
        callback()
    }, function () {
        swal("Error!", "Unknown error, please try again", "error")
    })
}

trial2.getDataPasivaPeriod = function(callback) {
    model.Processing(true)
    var url = "/financial/getdatapasivaperiod"
    var dateStart = $("#dateStart").data('kendoDatePicker').value();
    var dateEnd = $("#dateEnd").data("kendoDatePicker").value();
    var param = {
        DateStart: dateStart,
        DateEnd: dateEnd,
        Filter: trial.koFilterIsActive()
    }
    ajaxPost(url, param, function (res) {
        if (res.IsError === "true") {
            swal("Error!", res.Message, "error")
            return
        }
        var dataField = []
        _.forEach(res.Data.DataAcc, function (v) {
            dataField.push({
                beginingVal: v.Amount_Credit,
                creditVal: v.Amount_Credit,
                debetVal: v.Amount_Debet,
                endingVal: v.Amount_Ending
            })
        })
        var monthYear = []
        for (var d = new Date(dateStart);d <= new Date(dateEnd);d.setMonth(d.getMonth() + 1)) {
            var keyDate = moment(d).format("MMMYYYY")
            if (dataField.beginingVal == undefined) {
                dataField.beginingVal = 0
            }
            if (dataField.creditVal == undefined) {
                dataField.creditVal = 0
            }
            if (dataField.debetVal == undefined) {
                dataField.debetVal = 0
            }
            if (dataField.endingVal == undefined) {
                dataField.endingVal = 0
            }
            monthYear.push({
                M : parseInt(moment(d).format("M")),
                Y : parseInt(moment(d).format("YYYY")),
                MY : keyDate,
                my :  {
                         begining: dataField.beginingVal,
                         credit:dataField.creditVal,
                         debet: dataField.debetVal,
                         ending: dataField.endingVal
                   }
            })
        }
        var DataAmount = []
        _.forEach(res.Data.DataAmount, function(v,i) {
            var m = v.Month
            var y = v.Year
            var dateMy = new dateEnd(y,m-1, 01)
            var my = moment(dateMy).format("MMMYYYY")
            DataAmount.push({
                Acc_Code: v.Acc_Code,
                Acc_Name: v.Acc_Name,
                MonthYear: my,
                Amount: v.Amount
            })
        })
        var dataAmount2 = _.chain(DataAmount).groupBy("Acc_Code").toPairs().map(function (curentItem) {
            return _.zipObject(["Acc_Code", "DataItem"], curentItem);
        }).value();
        var dataAccP = []
        _.forEach(res.Data.DataAcc, function(v,i) {
            var name = v.Acc_Code
            var each = _.find(dataAmount2.DataItem, function(vv){return vv.MonthYear==name})
            if(each != undefined) {
                var itemMy = {}
                var totalAmount = 0
                _.forEach(monthYear, function(key) {
                    var eachItem = _.find(each.DataItem, function(vv){return vv.MonthYear==key.my})
                    if (eachItem == undefined) {
                        itemMy[key.MY] = 0
                    }else {
                        itemMy[key.my] = eachItem.Amount
                        totalAmount = totalAmount + eachItem
                    }
                })
                dataAccP.push({
                    Acc_Code: v.Acc_Code,
                    Main_Acc_Code: v.Main_Code,
                    Acc_Name: v.Acc_name,
                    DataItem: itemMy,
                    TotalAmount: totalAmount
                })
            }else {
                var itemMy = {}
                _.forEach(monthYear, function(key) {
                    itemMy[key.MY] = key.my
                })
                dataAccP.push({
                    Acc_Code: v.Acc_Code,
                    Main_Acc_Code: v.Main_Code,
                    Acc_Name: v.Acc_Name,
                    DataItem: itemMy,
                    TotalAmount: 0
                })
            }
        })
        trial2.dataPasiva(dataAccP)
        model.Processing(false)
        callback()
    }, function () {
        swal("Error!", "Unknown error, please try again", "error")
    })
}

trial2.getRenderPasivaPeriod = function() {
    var Data = trial2.dataPasiva()
    
    var columns = [{
        field: "Acc_Code",
        title: "Acc Code",
        width: 50
    }, {
        field: "Acc_Name",
        title: "Account Name",
        width: 70,
        footerTemplate: "<div stlye='text-align:center; font-size:15px;'>Total Pasiva:</div>"
    }]
    var dateStart = $("#dateStart").data("kendoDatePicker").value();
    var dateEnd = $("#dateEnd").data("kendoDatePicker").value();
    var Aggregate = []
    for(var d = new Date(dateStart); d <= new Date(dateEnd);d.setMonth(d.getMonth() + 1)){
        var namefield = "DataItem."+moment(d).format("MMMYYYY")+"."
        columns.push({
            field: namefield+"begining",
            title: moment(d).format("MMMYYYY")+'-Begining',
            // template: "#=ChangeToRupiah(sum)#",
            footerTemplate: "<div style='text-align:right;'>#=ChangeToRupiah(sum)#</div>",
            attributes: {
                style: "text-align:right;"
            },
            width: 100
        }, {
            field: namefield+"debet",
            title: moment(d).format("MMMYYYY")+'-Debet',
            // template: "#=ChangeToRupiah(sum)#",
            footerTemplate: "<div style='text-align:right;font-size:15px;'>#=ChangeToRupiah(sum)#</div>",
            attributes: {
                style: "text-align:right;"
            },
            width: 100
        }, {
            field: namefield+"credit",
            title: moment(d).format("MMMYYYY")+'-Credit',
            // template: "#=ChangeToRupiah(sum)#",
            width: 100,
            footerTemplate: "<div style='text-align:right;font-size:15px;'>#=ChangeToRupiah(sum)#</div>",
            attributes: {
                style: "text-align:right;"
            }
        }, {
            field: namefield+"ending",
            title: moment(d).format("MMMYYYY")+'-Ending',
            // template: "#=ChangeToRupiah(sum)#",
            width: 100,
            footerTemplate: "<div style='text-align:right;font-size:15px'>#=ChangeToRupiah(sum)#</div>",
            attributes: {
                style: "text-align:right;"
            }
        })
        Aggregate.push({
            field: namefield+"begining",
            aggregate: "sum"
        }, {
            field: namefield+"debet",
            aggregate: "sum"
        }, {
            field: namefield+"credit",
            aggregate: "sum"
        }, {
            field: namefield+"ending",
            aggregate: "sum"
        })
    }
    columns.push({
        field: "TotalAmount",
        title: "TOTAL",
        width: 100,
        attributes: {
            style: "text-align:right;"
        },
        hidden: true
    })
    $("#gridPasiva").html("")
    $("#gridPasiva").kendoGrid({
        dataSource: {
            data: Data,
            aggregate: Aggregate
        },
        scrollable: true,
        height: 400,
        columns: columns,
        columnmenu: true
    })
}

trial2.getDataIncomePeriod = function(callback) {
    var url = "/financial/getdatatrialincomeperiode"
    var dateStart = $("#dateStart").data("kendoDatePicker").value();
    var dateEnd = $("#dateEnd").data("kendoDatePicker").value();
    var param = {
        DateStart: dateStart,
        DateEnd: dateEnd,
        Filter: trial.koFilterIsActive(),
        Check: true
    }
    ajaxPost(url, param, function(res) {
        if (res.IsError === "true") {
            swal("Error!", res.Message, "error")
            return
        }
        var dataField = []
        _.forEach(res.Data.DataAcc, function (v) {
            dataField.push({
                beginingVal: v.Amount_Credit,
                creditVal: v.Amount_Credit,
                debetVal: v.Amount_Debet,
                endingVal: v.Amount_Ending
            })
        })
        var monthYear = []
        for (var d = new Date(dateStart);d <= new Date(dateEnd);d.setMonth(d.getMonth() + 1)) {
            var keyDate = moment(d).format("MMMYYYY")
            if (dataField.beginingVal == undefined) {
                dataField.beginingVal = 0
            }
            if (dataField.creditVal == undefined) {
                dataField.creditVal = 0
            }
            if (dataField.debetVal == undefined) {
                dataField.debetVal = 0
            }
            if (dataField.endingVal == undefined) {
                dataField.endingVal = 0
            }
            monthYear.push({
                M : parseInt(moment(d).format("M")),
                Y : parseInt(moment(d).format("YYYY")),
                MY : keyDate,
                my :  {
                         begining: dataField.beginingVal,
                         credit:dataField.creditVal,
                         debet: dataField.debetVal,
                         ending: dataField.endingVal
                   }
            })
        }
        var DataAmount = []
        _.forEach(res.Data.DataAmount, function(v,i) {
            var m = v.Month
            var y = v.Year
            var dateMy = new dateEnd(y,m-1, 01)
            var my = moment(dateMy).format("MMMYYYY")
            DataAmount.push({
                Acc_Code: v.Acc_Code,
                Acc_Name: v.Acc_Name,
                MonthYear: my,
                Amount: v.Amount
            })
        })
        var dataAmount2 = _.chain(DataAmount).groupBy("Acc_Code").toPairs().map(function (curentItem) {
            return _.zipObject(["Acc_Code", "DataItem"], curentItem);
        }).value();
        var dataAccInc = []
        _.forEach(res.Data.DataAcc, function(v,i) {
            var name = v.Acc_Code
            var each = _.find(dataAmount2.DataItem, function(vv){return vv.MonthYear==name})
            if(each != undefined) {
                var itemMy = {}
                var totalAmount = 0
                _.forEach(monthYear, function(key) {
                    var eachItem = _.find(each.DataItem, function(vv){return vv.MonthYear==key.my})
                    if (eachItem == undefined) {
                        itemMy[key.MY] = 0
                    }else {
                        itemMy[key.my] = eachItem.Amount
                        totalAmount = totalAmount + eachItem
                    }
                })
                dataAccInc.push({
                    Acc_Code: v.Acc_Code,
                    Main_Acc_Code: v.Main_Code,
                    Acc_Name: v.Acc_name,
                    DataItem: itemMy,
                    TotalAmount: totalAmount
                })
            }else {
                var itemMy = {}
                _.forEach(monthYear, function(key) {
                    itemMy[key.MY] = key.my
                })
                dataAccInc.push({
                    Acc_Code: v.Acc_Code,
                    Main_Acc_Code: v.Main_Code,
                    Acc_Name: v.Acc_Name,
                    DataItem: itemMy,
                    TotalAmount: 0
                })
            }
        })
        // console.log("dataAccPasiva :", dataAccP)
        trial2.dataIncomePeriode(dataAccInc)
        model.Processing(false)
        callback()
    }, function () {
        swal("Error!", "Unknown error, please try again", "error")
    })
}

trial2.rendergridIncome1 = function () {
    $('#gridIncomeStatement1').html("")
    //SALES AND REVENUE:5000
    var Data5001 = _.filter(trial2.dataIncomePeriode(), function(x){
        return x.Acc_Code == 5110
    });
    for (i in Data5001) {
        Data5001[i].Parent = "PENJUALAN"
    }
    var data = Data5001
    var columns = [{
        field: 'Acc_Code',
        title: 'Acc Code',
        width: 50
    }, {
        field: 'Acc_Name',
        title: 'Account Name',
        width: 150,
        footerTemplate: "<div style='text-align:left; font-size: 15px;'>TOTAL PENJUALAN:</div>"
    }, {
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #",
        hidden: true
    }]
    var dateStart = $("#dateStart").data("kendoDatePicker").value();
    var dateEnd = $("#dateEnd").data("kendoDatePicker").value();
    var Aggregate = []
    for (var d = new Date(dateStart); d <= new Date(dateEnd); d.setMonth(d.getMonth() + 1)) {
        var namefield = "DataItem." +moment(d).format("MMMYYYY")+"."
        columns.push(
        //     title: moment(d).format("MMM YYYY"),
        //     field: namefield,
        //     template: "#=ChangeToRupiah("+namefield+")#",
        //     width: 100,
        //     attributes: {
        //         style: "text-align:right;"
        //     },
        //     // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",#= ChangeToRupiah(income2.sumSalesEndingLabaKotorPeriod('"+namefield+"', sum)) #, #= ChangeToRupiah(income2.sumTotalHPP("+namefield+")) #
        //     footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income2.sumSalesEndingLabaKotorPeriod('"+namefield+"', sum)) #</div><div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income2.sumTotalHPP('"+namefield+"')) #</div>",
        // }
        {
        field: namefield+'begining',
        title: moment(d).format("MMMYYYY")+'-Begining',
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: namefield+'debet',
        title: moment(d).format("MMMYYYY")+'-Debet',
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: namefield+'credit',
        title: moment(d).format("MMMYYYY")+'-Credit',
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: namefield+'ending',
        title: moment(d).format("MMMYYYY")+'-Ending',
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    })
    Aggregate.push({
        field: namefield+"begining",
        aggregate: "sum"
    }, {
        field: namefield+"debet",
        aggregate: "sum"
    },{
        field: namefield+"credit",
        aggregate: "sum"
    }, {
        field: namefield+"ending",
        aggregate: "sum"
    })
}
    $("#gridIncomeStatement1").kendoGrid({
        dataSource: {
            data: data,
            aggregate: Aggregate,
            group: {
                dir: "asc",
                field: "Parent"
            }
        },
        scrollable: true,
        columns: columns,
        columnmenu: true
    })
    // var grid = $("#gridIncomeStatement1").data("kendoGrid");
    // $("#gridIncomeSales .k-grid-footer").css('display', 'none');
}

trial2.rendergridIncome1Retur = function () {
    $('#gridIncomeStatement1retur').html("")
    var dataretur = _.filter(trial2.dataIncomePeriode(), function(x){
        return x.Acc_Code== 5210 || x.Acc_Code==5211 || x.Acc_Code==5212  
    });
    for (i in dataretur) {
        dataretur[i].Parent = "HARGA POKOK"
    }
    var data = dataretur
    var columns = [{
        field: 'Acc_Code',
        title: 'Acc Code',
        width: 50
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
    var dateStart = $("#dateStart").data("kendoDatePicker").value();
    var dateEnd = $("#dateEnd").data("kendoDatePicker").value();
    var Aggregate = []
    for (var d = new Date(dateStart); d <= new Date(dateEnd); d.setMonth(d.getMonth() + 1)) {
        var namefield = "DataItem." +moment(d).format("MMMYYYY")+"."
        columns.push(
        //     title: moment(d).format("MMM YYYY"),
        //     field: namefield,
        //     template: "#=ChangeToRupiah("+namefield+")#",
        //     width: 100,
        //     attributes: {
        //         style: "text-align:right;"
        //     },
        //     // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",#= ChangeToRupiah(income2.sumSalesEndingLabaKotorPeriod('"+namefield+"', sum)) #, #= ChangeToRupiah(income2.sumTotalHPP("+namefield+")) #
        //     footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income2.sumSalesEndingLabaKotorPeriod('"+namefield+"', sum)) #</div><div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income2.sumTotalHPP('"+namefield+"')) #</div>",
        // }
        {
        field: namefield+'begining',
        title: moment(d).format("MMMYYYY")+'-Begining',
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: namefield+'debet',
        title: moment(d).format("MMMYYYY")+'-Debet',
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: namefield+'credit',
        title: moment(d).format("MMMYYYY")+'-Credit',
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: namefield+'ending',
        title: moment(d).format("MMMYYYY")+'-Ending',
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    })
    Aggregate.push({
        field: namefield+"begining",
        aggregate: "sum"
    }, {
        field: namefield+"debet",
        aggregate: "sum"
    },{
        field: namefield+"credit",
        aggregate: "sum"
    }, {
        field: namefield+"ending",
        aggregate: "sum"
    })
}
    $("#gridIncomeStatement1retur").kendoGrid({
        dataSource: {
            data: data,
            aggregate: Aggregate,
            group: {
                dir: "asc",
                field: "Parent"
            }
        },
        scrollable: true,
        columns: columns,
        columnmenu: true
    })
}

trial2.rendergridIncome2 = function () {
    $('#gridIncomeStatement2').html("")
    //Operating Expenses:6000
    var Data6000 = _.filter(trial2.dataIncomePeriode(), function(x){
        return x.Acc_Code!=6160&&x.Acc_Code!=6161&&x.Acc_Code!=6999&& x.Main_Acc_Code == 6000
     });
    for (i in Data6000) {
        Data6000[i].Parent = "BIAYA UMUM DAN ADMINISTRASI"
    }
    var data = Data6000

    var columns = [{
        field: 'Acc_Code',
        title: 'Acc Code',
        width: 50
    }, {
        field: 'Acc_Name',
        title: 'Account Name',
        width: 150,
        footerTemplate: "<div style='text-align:left; font-size: 15px;'>TOTAL BIAYA UMUM DAN ADMINISTRASI:</div><div style='text-align:left; font-size: 15px;'>LABA USAHA:</div>"
    }, {
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #",
        hidden: true,
    }]
    var dateStart = $("#dateStart").data("kendoDatePicker").value();
    var dateEnd = $("#dateEnd").data("kendoDatePicker").value();
    var Aggregate = []
    for (var d = new Date(dateStart); d <= new Date(dateEnd); d.setMonth(d.getMonth() + 1)) {
        var namefield = "DataItem." +moment(d).format("MMMYYYY")+"."
        columns.push(
        //     title: moment(d).format("MMM YYYY"),
        //     field: namefield,
        //     template: "#=ChangeToRupiah("+namefield+")#",
        //     width: 100,
        //     attributes: {
        //         style: "text-align:right;"
        //     },
        //     // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",#= ChangeToRupiah(income2.sumSalesEndingLabaKotorPeriod('"+namefield+"', sum)) #, #= ChangeToRupiah(income2.sumTotalHPP("+namefield+")) #
        //     footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income2.sumSalesEndingLabaKotorPeriod('"+namefield+"', sum)) #</div><div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income2.sumTotalHPP('"+namefield+"')) #</div>",
        // }
        {
        field: namefield+'begining',
        title: moment(d).format("MMMYYYY")+'-Begining',
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: namefield+'debet',
        title: moment(d).format("MMMYYYY")+'-Debet',
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: namefield+'credit',
        title: moment(d).format("MMMYYYY")+'-Credit',
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: namefield+'ending',
        title: moment(d).format("MMMYYYY")+'-Ending',
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    })
    Aggregate.push({
        field: namefield+"begining",
        aggregate: "sum"
    }, {
        field: namefield+"debet",
        aggregate: "sum"
    },{
        field: namefield+"credit",
        aggregate: "sum"
    }, {
        field: namefield+"ending",
        aggregate: "sum"
    })
}
    $("#gridIncomeStatement2").kendoGrid({
        dataSource: {
            data: data,
            aggregate: Aggregate,
            group: {
                dir: "asc",
                field: "Parent"
            }
        },
        scrollable: true,
        columns: columns,
        columnmenu: true
    })
}

trial2.rendergridIncome3 = function () {
    $("#gridIncomeStatement3").html("")
    var Data7000 = _.filter(trial2.dataIncomePeriode(), function (x){
        return x.Main_Acc_Code == 7000
    });
    for (i in Data7000) {
        Data7000[i].Parent = "PENDAPATAN/(BIAYA) DI LUAR USAHA"
    }
    var data = Data7000

    var columns = [{
        field: 'Acc_Code',
        title: 'Acc Code',
        width: 50
    }, {
        field: 'Acc_Name',
        title: 'Account Name',
        width: 150,
        footerTemplate: "<div style='text-align:left; font-size: 15px;'>TOTAL PENDAPATAN/(BIAYA) DI LUAR USAHA:</div>"
    }, {
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #",
        hidden: true,
    }]
    var dateStart = $("#dateStart").data("kendoDatePicker").value();
    var dateEnd = $("#dateEnd").data("kendoDatePicker").value();
    var Aggregate = []
    for (var d = new Date(dateStart); d <= new Date(dateEnd); d.setMonth(d.getMonth() + 1)) {
        var namefield = "DataItem." +moment(d).format("MMMYYYY")+"."
        columns.push(
        //     title: moment(d).format("MMM YYYY"),
        //     field: namefield,
        //     template: "#=ChangeToRupiah("+namefield+")#",
        //     width: 100,
        //     attributes: {
        //         style: "text-align:right;"
        //     },
        //     // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",#= ChangeToRupiah(income2.sumSalesEndingLabaKotorPeriod('"+namefield+"', sum)) #, #= ChangeToRupiah(income2.sumTotalHPP("+namefield+")) #
        //     footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income2.sumSalesEndingLabaKotorPeriod('"+namefield+"', sum)) #</div><div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income2.sumTotalHPP('"+namefield+"')) #</div>",
        // }
        {
        field: namefield+'begining',
        title: moment(d).format("MMMYYYY")+'-Begining',
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: namefield+'debet',
        title: moment(d).format("MMMYYYY")+'-Debet',
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: namefield+'credit',
        title: moment(d).format("MMMYYYY")+'-Credit',
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: namefield+'ending',
        title: moment(d).format("MMMYYYY")+'-Ending',
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    })
    Aggregate.push({
        field: namefield+"begining",
        aggregate: "sum"
    }, {
        field: namefield+"debet",
        aggregate: "sum"
    },{
        field: namefield+"credit",
        aggregate: "sum"
    }, {
        field: namefield+"ending",
        aggregate: "sum"
    })
}
    $("#gridIncomeStatement3").kendoGrid({
        dataSource: {
            data: data,
            aggregate: Aggregate,
            group: {
                dir: "asc",
                field: "Parent"
            }
        },
        scrollable: true,
        columns: columns,
        columnmenu: true
    })
}

trial2.rendergridIncome4 = function() {
    $("#gridIncomeStatement4").html("")
    var Data8000 = _.filter(trial2.dataIncomePeriode(), {
        Main_Acc_Code: 8000
    });
    for (i in Data8000) {
        Data8000[i].Parent = "BIAYA DI LUAR USAHA"
    }
    var data = Data8000

    var columns = [{
        field: 'Acc_Code',
        title: 'Acc Code',
        width: 50
    }, {
        field: 'Acc_Name',
        title: 'Account Name',
        width: 150,
        footerTemplate: "<div style='text-align:left; font-size: 15px;'>TOTAL PENDAPATAN/(BIAYA) DI LUAR USAHA:</div>"
    }, {
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #",
        hidden: true,
    }]
    var dateStart = $("#dateStart").data("kendoDatePicker").value();
    var dateEnd = $("#dateEnd").data("kendoDatePicker").value();
    var Aggregate = []
    for (var d = new Date(dateStart); d <= new Date(dateEnd); d.setMonth(d.getMonth() + 1)) {
        var namefield = "DataItem." +moment(d).format("MMMYYYY")+"."
        columns.push(
        //     title: moment(d).format("MMM YYYY"),
        //     field: namefield,
        //     template: "#=ChangeToRupiah("+namefield+")#",
        //     width: 100,
        //     attributes: {
        //         style: "text-align:right;"
        //     },
        //     // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",#= ChangeToRupiah(income2.sumSalesEndingLabaKotorPeriod('"+namefield+"', sum)) #, #= ChangeToRupiah(income2.sumTotalHPP("+namefield+")) #
        //     footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income2.sumSalesEndingLabaKotorPeriod('"+namefield+"', sum)) #</div><div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income2.sumTotalHPP('"+namefield+"')) #</div>",
        // }
        {
        field: namefield+'begining',
        title: moment(d).format("MMMYYYY")+'-Begining',
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: namefield+'debet',
        title: moment(d).format("MMMYYYY")+'-Debet',
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: namefield+'credit',
        title: moment(d).format("MMMYYYY")+'-Credit',
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: namefield+'ending',
        title: moment(d).format("MMMYYYY")+'-Ending',
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    })
    Aggregate.push({
        field: namefield+"begining",
        aggregate: "sum"
    }, {
        field: namefield+"debet",
        aggregate: "sum"
    },{
        field: namefield+"credit",
        aggregate: "sum"
    }, {
        field: namefield+"ending",
        aggregate: "sum"
    })
}
    $("#gridIncomeStatement4").kendoGrid({
        dataSource: {
            data: data,
            aggregate: Aggregate,
            group: {
                dir: "asc",
                field: "Parent"
            }
        },
        scrollable: true,
        columns: columns,
        columnmenu: true
    })
}

trial2.getDataTaxPeriode = function (callback) {
    model.Processing(true)
    var url = "/financial/getdatataxperiod"
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var param = {}
    if (trial.koFilterIsActive() == true) {
        param = {
            DateStart: dateStart,
            DateEnd: dateEnd,
            TextSearch: trial.textSearch(),
            Filter: true
        }
    }else {
        param = {
            Filter: false,
        }
    }

    ajaxPost(url, param, function (res) {
        if (res.IsError === "true") {
            swal("Error!", res.Message, "error")
            return
        }
        var dataField = []
        _.forEach(res.Data.DataAcc, function (v) {
            dataField.push({
                beginingVal: v.Amount_Credit,
                creditVal: v.Amount_Credit,
                debetVal: v.Amount_Debet,
                endingVal: v.Amount_Ending
            })
        })
        var monthYear = []
        for (var d = new Date(dateStart);d <= new Date(dateEnd);d.setMonth(d.getMonth() + 1)) {
            var keyDate = moment(d).format("MMMYYYY")
            if (dataField.beginingVal == undefined) {
                dataField.beginingVal = 0
            }
            if (dataField.creditVal == undefined) {
                dataField.creditVal = 0
            }
            if (dataField.debetVal == undefined) {
                dataField.debetVal = 0
            }
            if (dataField.endingVal == undefined) {
                dataField.endingVal = 0
            }
            monthYear.push({
                M : parseInt(moment(d).format("M")),
                Y : parseInt(moment(d).format("YYYY")),
                MY : keyDate,
                my :  {
                         begining: dataField.beginingVal,
                         credit:dataField.creditVal,
                         debet: dataField.debetVal,
                         ending: dataField.endingVal
                   }
            })
        }
        var DataAmount = []
        _.forEach(res.Data.DataAmount, function(v,i) {
            var m = v.Month
            var y = v.Year
            var dateMy = new dateEnd(y,m-1, 01)
            var my = moment(dateMy).format("MMMYYYY")
            DataAmount.push({
                Acc_Code: v.Acc_Code,
                Acc_Name: v.Acc_Name,
                MonthYear: my,
                Amount: v.Amount
            })
        })
        var dataAmount2 = _.chain(DataAmount).groupBy("Acc_Code").toPairs().map(function (curentItem) {
            return _.zipObject(["Acc_Code", "DataItem"], curentItem);
        }).value();
        var dataAcc = []
        _.forEach(res.Data.DataAcc, function(v,i) {
            var name = v.Acc_Code
            var each = _.find(dataAmount2.DataItem, function(vv){return vv.MonthYear==name})
            if(each != undefined) {
                var itemMy = {}
                var totalAmount = 0
                _.forEach(monthYear, function(key) {
                    var eachItem = _.find(each.DataItem, function(vv){return vv.MonthYear==key.my})
                    if (eachItem == undefined) {
                        itemMy[key.MY] = 0
                    }else {
                        itemMy[key.my] = eachItem.Amount
                        totalAmount = totalAmount + eachItem
                    }
                })
                dataAcc.push({
                    Acc_Code: v.Acc_Code,
                    Main_Acc_Code: v.Main_Code,
                    Acc_Name: v.Acc_name,
                    DataItem: itemMy,
                    TotalAmount: totalAmount
                })
            }else {
                var itemMy = {}
                _.forEach(monthYear, function(key) {
                    itemMy[key.MY] = key.my
                })
                dataAcc.push({
                    Acc_Code: v.Acc_Code,
                    Main_Acc_Code: v.Main_Code,
                    Acc_Name: v.Acc_Name,
                    DataItem: itemMy,
                    TotalAmount: 0
                })
            }
        })
        // console.log("dataAcc :", dataAcc)
        trial2.dataTax(dataAcc)
        model.Processing(false)
        callback()
    }, function () {
        swal("Error!", "Unknown error, please try again", "error")
    })
}

trial2.renderGridTax = function () {
    $('#gridTAX').html("")
    var dataPeriod = trial2.dataTax()
    for (i in dataPeriod) {
        dataPeriod[i].Parent = "PAJAK"
    }

    // console.log("data =>", data)
    var columns = [{
        field: 'Acc_Code',
        title: 'Acc Code',
        width: 50
    }, {
        field: 'Acc_Name',
        title: 'Account Name',
        width: 150,
        footerTemplate: "<div style='text-align:left; font-size: 15px;'>EARNING AFTER TAX:</div>"
    }, {
        field: 'Parent',
        title: 'Parent',
        width: 50,
        groupHeaderTemplate: "#= value #",
        hidden: true,
    }]
    var dateStart = $("#dateStart").data("kendoDatePicker").value();
    var dateEnd = $("#dateEnd").data("kendoDatePicker").value();
    var Aggregate = []
    for (var d = new Date(dateStart); d <= new Date(dateEnd); d.setMonth(d.getMonth() + 1)) {
        var namefield = "DataItem." +moment(d).format("MMMYYYY")+"."
        var dy = moment(d).format("MMMYYYY")
        columns.push(
        //     title: moment(d).format("MMM YYYY"),
        //     field: namefield,
        //     template: "#=ChangeToRupiah("+namefield+")#",
        //     width: 100,
        //     attributes: {
        //         style: "text-align:right;"
        //     },
        //     // footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",#= ChangeToRupiah(income2.sumSalesEndingLabaKotorPeriod('"+namefield+"', sum)) #, #= ChangeToRupiah(income2.sumTotalHPP("+namefield+")) #
        //     footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income2.sumSalesEndingLabaKotorPeriod('"+namefield+"', sum)) #</div><div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(income2.sumTotalHPP('"+namefield+"')) #</div>",
        // }
        {
        field: namefield+'begining',
        title: dy+"-Begining",
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: namefield+'debet',
        title: dy+"-Debet",
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: namefield+'credit',
        title: dy+"-Credit",
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: namefield+'ending',
        title: dy+"-Ending",
        // template: "#=ChangeToRupiah(sum)#",
        width: 50,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    })
    Aggregate.push({
        field: namefield+"begining",
        aggregate: "sum"
    }, {
        field: namefield+"debet",
        aggregate: "sum"
    },{
        field: namefield+"credit",
        aggregate: "sum"
    }, {
        field: namefield+"ending",
        aggregate: "sum"
    })
}
    $("#gridTAX").kendoGrid({
        dataSource: {
            data: dataPeriod,
            aggregate: Aggregate,
            group: {
                dir: "asc",
                field: "Parent"
            }
        },
        scrollable: true,
        columns: columns,
        columnmenu: true
    })
}