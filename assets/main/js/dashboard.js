var dashboard = {}
dashboard.koDataChartRev = ko.observableArray([]);
dashboard.koDataTotRev = ko.observable();
dashboard.koDataPrevRev = ko.observable();
dashboard.koDataPercentRev = ko.observable();
dashboard.koDataChartExp = ko.observableArray([]);
dashboard.koDataTotExp = ko.observable();
dashboard.koDataPrevExp = ko.observable();
dashboard.koDataPercentExp = ko.observable();
dashboard.koDataChartNet = ko.observableArray([]);
dashboard.koDataTotNet = ko.observable();
dashboard.koDataPrevNet = ko.observable();
dashboard.koDataPercentNet = ko.observable();
dashboard.koDateNow = ko.observable("");
dashboard.koDatePrev = ko.observable();
dashboard.koDataChartRevExNet = ko.observableArray([]);
dashboard.koDataTotRevExNet = ko.observable();
dashboard.koDataMonthlyRevenue = ko.observable();
dashboard.koDataMonthlyExpenses = ko.observable();
dashboard.koDataMonthlyRevenueVal = ko.observable();
dashboard.koDataMonthlyRevPrev = ko.observable();
dashboard.koDataMonthlyExpensesVal = ko.observable();
dashboard.koDataMonthlyExpPrev = ko.observable();
dashboard.koDateMonthlyNow = ko.observable();
dashboard.koDateMonthlyPrev = ko.observable();
dashboard.koDatePageBar = ko.observable();
dashboard.koMonthYear = ko.observable();
dashboard.koData5topExpenses = ko.observableArray([])
dashboard.koData5topRevenue = ko.observableArray([])
dashboard.koDataCurrentAssets = ko.observableArray([])
dashboard.koDataINV = ko.observableArray([])
dashboard.koIDRamount = ko.observable()
dashboard.koUSDamount = ko.observable()
dashboard.koPTCamount = ko.observable()
dashboard.koPETTYamount = ko.observable()
dashboard.koTotalOSInvoice = ko.observable()
dashboard.koTotalOSPO = ko.observable()
dashboard.koLenOSInvoice = ko.observable()
dashboard.koLenOSPO = ko.observable()
dashboard.titleDashboard = ko.observable("Dashboard Summary")
dashboard.visibleDashboard = ko.observable(false)
model.Processing(true)
dashboard.getDateNow = function () {
    var start = moment().startOf('year').format("DD MMMM YYYY")
    var now = moment().startOf('day').format("DD MMMM YYYY")
    var range = start + " to " + now

    var dateFrom = moment(dateFrom).startOf('year').subtract(1, 'year').format("DD MMMM YYYY");
    var dateEnd = moment(dateFrom).endOf('year').subtract(dateFrom, 'year').format("DD MMMM YYYY");
    var rangePrev = dateFrom + " to " + dateEnd + " (Prev)"
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    dashboard.koDateNow(range)
    dashboard.koDatePrev(rangePrev)
    dashboard.koDatePageBar(page)
}

dashboard.getDateMonthly = function () {
    var start = moment().startOf('month').format("DD MMMM YYYY")
    var now = moment().startOf('day').format("DD MMMM YYYY")
    var range = start + " to " + now
    var monthyear = moment().startOf('month').format("MMMM YYYY")

    var dateFrom = moment(dateFrom).startOf('month').subtract(1, 'month').format("DD MMMM YYYY");
    var dateEnd = moment(dateFrom).endOf('month').subtract(dateFrom, 'month').format("DD MMMM YYYY");
    var rangePrev = dateFrom + " to " + dateEnd + " (Prev)"
    dashboard.koDateMonthlyNow(range)
    dashboard.koDateMonthlyPrev(rangePrev)
    dashboard.koMonthYear(monthyear)
}

dashboard.getDataChartRevenue = function (callback) {
    var url = "/dashboard/getdatachartrevenue"
    var param = {
        Year: parseInt(moment().format("YYYY"))
    }
    model.Processing(true)
    ajaxPost(url, param, function (res) {
        var Data = res.Data.Chart
        Data.push({
            _id: 0,
            Revenue: 0,
        })
        Data = _.sortBy(Data, [function (o) {
            return o._id;
        }]);
        var total = res.Data.Total
        if (total < 0) {
            dashboard.koDataTotRev("(" + kendo.toString(total * -1, "n2") + ")")
        } else {
            dashboard.koDataTotRev(kendo.toString(total, "n2"))
        }

        var PrevTotal = res.Data.PrevTotal
        dashboard.koDataPrevRev(kendo.toString(PrevTotal, "n2"))
        if (total >= PrevTotal) {
            var a = total - PrevTotal
            if (total == 0) {
                dashboard.koDataPercentRev(100 + "%")
            } else {
                dashboard.koDataPercentRev(((a / total) * 100).toFixed(0) + "%")
            }
            document.getElementById("RevenueI").style.color = "green";
        } else {
            var b = PrevTotal - total
            if (PrevTotal == 0) {
                dashboard.koDataPercentRev(100 + "%")
            } else if (PrevTotal > 0) {
                dashboard.koDataPercentRev(((b / PrevTotal) * 100).toFixed(0) + "%")
            } else if (PrevTotal < 0) {
                dashboard.koDataPercentRev(((b / total) * 100).toFixed(0) + "%")
            }

            document.getElementById("RevenueI").style.color = "red";
            $('#RevenueI').removeClass('fa fa-caret-up font').addClass('fa fa-caret-down font');

        }
        dashboard.koDataChartRev(Data)
        callback()
        model.Processing(false)
    })
}

dashboard.getDataChartExpenses = function (callback) {
    var url = "/dashboard/getdatachartexpenses"
    var param = {
        Year: parseInt(moment().format("YYYY"))
    }
    model.Processing(true)
    ajaxPost(url, param, function (res) {
        var Data = res.Data.Chart
        Data.push({
            _id: 0,
            Expenses: 0,
        })

        Data = _.sortBy(Data, [function (o) {
            return o._id;
        }]);

        var total = res.Data.Total
        if (total < 0) {
            dashboard.koDataTotExp("(" + kendo.toString(total * -1, "n2") + ")")
        } else {
            dashboard.koDataTotExp(kendo.toString(total, "n2"))
        }

        var PrevTotal = res.Data.PrevTotal

        dashboard.koDataPrevExp(kendo.toString(PrevTotal, "n2"))
        if (total >= PrevTotal) {
            var a = total - PrevTotal
            if (total == 0) {
                dashboard.koDataPercentExp(100 + "%")
            } else {
                dashboard.koDataPercentExp(((a / total) * 100).toFixed(0) + "%")
            }

            document.getElementById("ExpensesI").style.color = "green";
        } else {
            var b = PrevTotal - total
            if (PrevTotal == 0) {
                dashboard.koDataPercentExp(100 + "%")
            } else if (PrevTotal > 0) {
                dashboard.koDataPercentExp(((b / PrevTotal) * 100).toFixed(0) + "%")
            } else if (PrevTotal < 0) {
                dashboard.koDataPercentExp(((b / total) * 100).toFixed(0) + "%")
            }

            document.getElementById("ExpensesI").style.color = "red";
            $('#ExpensesI').removeClass('fa fa-caret-up font').addClass('fa fa-caret-down font');

        }

        dashboard.koDataChartExp(Data)
        callback()
        model.Processing(false)
    })
}

dashboard.getDataChartNetProfit = function (callback) {
    var url = "/dashboard/getdatachartnetprofit"
    var param = {
        Year: parseInt(moment().format("YYYY"))
    }
    model.Processing(true)
    ajaxPost(url, param, function (res) {
        var Data = res.Data.Chart
        Data.push({
            _id: 0,
            Amount: 0,
        })
        Data = _.sortBy(Data, [function (o) {
            return o._id;
        }]);

        var total = res.Data.Total
        if (total < 0) {
            dashboard.koDataTotNet("(" + kendo.toString(total * -1, "n2") + ")")
        } else {
            dashboard.koDataTotNet(kendo.toString(total, "n2"))
        }

        var PrevTotal = res.Data.PrevTotal
        dashboard.koDataPrevNet(kendo.toString(PrevTotal, "n2"))
        if (total >= PrevTotal) {
            var a = total - PrevTotal
            if (total == 0) {
                dashboard.koDataPercentNet(100 + "%")
            } else {
                dashboard.koDataPercentNet(((a / total) * 100).toFixed(0) + "%")
            }

            document.getElementById("NetProfitI").style.color = "green";
        } else {
            var b = PrevTotal - total
            if (PrevTotal == 0) {
                dashboard.koDataPercentNet(100 + "%")
            } else if (PrevTotal > 0) {
                dashboard.koDataPercentNet(((b / PrevTotal) * 100).toFixed(0) + "%")
            } else if (PrevTotal < 0) {
                dashboard.koDataPercentNet(((b / total) * 100).toFixed(0) + "%")
            }

            document.getElementById("NetProfitI").style.color = "red";
            $('#NetProfitI').removeClass('fa fa-caret-up font').addClass('fa fa-caret-down font');

        }

        dashboard.koDataChartNet(Data)
        callback()
        model.Processing(false)
    })
}

dashboard.getDataChartRevExNet = function (callback) {
    var url = "/dashboard/getdatachartrevexnet"
    var param = {
        Year: parseInt(moment().format("YYYY"))
    }
    model.Processing(true)
    ajaxPost(url, param, function (res) {
        var Data = res.Data
        Data = _.sortBy(Data, [function (o) {
            return o._id;
        }]);
        for (i in Data) {
            Data[i].Date = new Date(moment().format("YYYY") + "-" + Data[i].Month + "-" + "1")
        }
        dashboard.koDataChartRevExNet(Data)

        callback()
        model.Processing(false)
    })
}
dashboard.renderChartRev = function () {


    $("#chartRev").kendoChart({
        dataSource: dashboard.koDataChartRev(),
        seriesColors: ["#4c07e0", "#4c07e0"],
        seriesDefaults: {
            type: "area"
        },
        legend: {
            visible: false,
        },
        series: [{
            field: "Revenue",
            name: "Revenue",
            line: {
                color: '#321177',
                width: 2
            }
        }],
        categoryAxis: {
            visible: false,
            field: "_id",
            majorGridLines: {
                visible: false
            },
            minorGridLines: {
                visible: false
            },
        },
        valueAxis: {
            visible: false,
            majorGridLines: {
                visible: false
            },
            minorGridLines: {
                visible: false
            },
        },
        tooltip: {
            visible: true,
            template: function (e) {
                var rev = e.dataItem.Revenue
                if (rev < 0) {
                    return "(" + kendo.toString(rev * -1, "n2") + ")"
                } else {
                    return kendo.toString(rev, "n2")
                }
            }
        },
        chartArea: {
            height: 60
        },
    });


}
dashboard.renderChartExp = function () {
    $("#chart1").kendoChart({
        dataSource: dashboard.koDataChartExp(),
        seriesColors: ["#cb0ed3", "#cb0ed3"],
        seriesDefaults: {
            type: "area"
        },
        legend: {
            visible: false,
        },
        series: [{
            field: "Expenses",
            name: "Expenses",
            line: {
                color: '#69116d',
                width: 2
            }
        }],
        categoryAxis: {
            visible: false,
            field: "_id",
            majorGridLines: {
                visible: false
            },
            minorGridLines: {
                visible: false
            },
        },
        valueAxis: {
            visible: false,
            majorGridLines: {
                visible: false
            },
            minorGridLines: {
                visible: false
            },
        },
        tooltip: {
            visible: true,
            template: function (e) {
                var exp = e.dataItem.Expenses
                if (exp < 0) {
                    return "(" + kendo.toString(exp * -1, "n2") + ")"
                } else {
                    return kendo.toString(exp, "n2")
                }
            }
        },
        chartArea: {
            height: 60
        },
    });
}
dashboard.renderChartNet = function () {
    $("#chart2").kendoChart({
        dataSource: {
            data: dashboard.koDataChartNet(),
            sort: [{
                field: "Month",
                dir: "asc"
            }],
        },
        seriesColors: ["#e20bb0", "#e20bb0"],
        seriesDefaults: {
            type: "area"
        },
        legend: {
            visible: false,
        },
        series: [{
            field: "Amount",
            name: "Amount",
            line: {
                color: '#af0e69',
                width: 2
            }
        }],
        categoryAxis: {
            visible: false,
            field: "Month",
            majorGridLines: {
                visible: false
            },
            minorGridLines: {
                visible: false
            },
        },

        valueAxis: {
            visible: false,
            majorGridLines: {
                visible: false
            },
            minorGridLines: {
                visible: false
            },
        },
        tooltip: {
            visible: true,
            template: function (e) {
                var net = e.dataItem.Amount
                if (net < 0) {
                    return "(" + kendo.toString(net * -1, "n2") + ")"
                } else {
                    return kendo.toString(net, "n2")
                }
            }
        },
        chartArea: {
            height: 60
        },
    });
}
dashboard.renderChartRevExNet = function () {
    var blueColor = ["#8bc6ed", "#5497c4", "#2e6c96", "#1b4f72", "#1072b5", "#0b8ce5", "#31c5f7", "#042c44", "#023d60", "#00507f"]
    var Data = new kendo.data.DataSource({
        data: dashboard.koDataChartRevExNet(),
        group: {
            field: "Type"
        },
        sort: [{
            field: "Date",
            dir: "asc"
        }, {
            field: "Type",
            dir: "asc"
        }],
        schema: {
            model: {
                fields: {
                    Date: {
                        type: "date"
                    },
                    Type: {
                        type: "string"
                    }
                }
            }
        }
    });

    $("#coulumn").kendoChart({
        dataSource: Data,
        legend: {
            position: "top",

        },
        seriesDefaults: {
            type: "column"
        },
        title: {
            text: dashboard.koDateNow(),
            font: "bold 15px Arial,Helvetica,sans-serif",
            color: "black"
        },
        seriesColors: ["#8cc739", "#085dad", "#f4222d"],
        series: [{
            name: "#= group.value #",
            field: "Amount",
            overlay: {
                gradient: "none"
            },
        }],
        valueAxis: {
            majorGridLines: {
                visible: false
            },
            minorGridLines: {
                visible: false
            },
            labels: {
                template: "#= dashboard.shortLabels(value) #"
            }
        },
        categoryAxis: {
            majorGridLines: {
                visible: false
            },
            minorGridLines: {
                visible: false
            },
            field: "Date",
            labels: {
                font: "11px sans-serif",
                format: "MMM\nyyyy"
            }
        },
        tooltip: {
            visible: true,
            template: "#= series.name #: #= ChangeToRupiah(value)#"
        },
        chartArea: {
            height: 353,
        },
    });
}
dashboard.getDataChartMonthlyRevenue = function (callback) {
    var url = "/dashboard/monthlyrevenue"
    model.Processing(true)
    ajaxPost(url, {}, function (res) {
        model.Processing(false)
        dashboard.koDataMonthlyRevenue(res.Data)
        var Data = res.Data[0]
        if (Data < 0) {
            dashboard.koDataMonthlyRevenueVal("(" + kendo.toString(Data * -1, "n2") + ")")
        } else {
            dashboard.koDataMonthlyRevenueVal(kendo.toString(Data, "n2"))
        }

        var DataPrev = res.Data[1]
        if (DataPrev < 0) {
            dashboard.koDataMonthlyRevPrev("(" + kendo.toString(DataPrev * -1, "n2") + ")")
        } else {
            dashboard.koDataMonthlyRevPrev(kendo.toString(DataPrev, "n2") + " (Prev Month)")
        }


        callback();
    });
}
dashboard.renderChartMonthlyRevenue = function () {
    data = dashboard.koDataMonthlyRevenue()
    var max = 0
    var biValue = 0
    if (data[0] > data[1]) {
        var value = 10 / 100 * data[0]
        max = data[0] + value
        biValue = data[0]
    } else {
        var value = 10 / 100 * data[1]
        max = data[1] + value
        biValue = data[1]
    }
    $("#monthlyRevenue").kendoChart({
        seriesColors: ["#f4222d", "#195781"],
        seriesDefaults: {
            type: "bullet"
        },
        series: [{
            target: {
                border: {
                    width: 3,
                    color: "#195781"
                }
            },
            data: [dashboard.koDataMonthlyRevenue()],
            labels: {
                visible: true,
                background: "transparent",
                template: "#= kendo.toString(value, n2) #"
            }
        }],
        categoryAxis: {
            majorGridLines: {
                visible: false
            },
            majorTicks: {
                visible: false
            }
        },
        valueAxis: {
            plotBands: [{
                from: 0,
                to: data[1],
                color: "#3498db",
                opacity: .6
            }],
            visible: false,
            majorGridLines: {
                visible: false
            },
            minorGridLines: {
                visible: false
            },
            minorTicks: {
                visible: true
            },
            min: 0,
            max: max,
        },

        tooltip: {
            visible: true,
            template: "Previous: #= kendo.toString(value.target, 'n2') # <br /> Current: #= kendo.toString(value.current, 'n2') #"
        },
        chartArea: {
            height: 60
        },
    });
}
dashboard.getDataChartMonthlyExpenses = function (callback) {
    var url = "/dashboard/monthlyexpenses"
    model.Processing(true)
    ajaxPost(url, {}, function (res) {
        model.Processing(false)
        dashboard.koDataMonthlyRevenue(res.Data)
        var Data = res.Data[0]
        if (Data < 0) {
            dashboard.koDataMonthlyExpensesVal("(" + kendo.toString(Data * -1, "n2") + ")")
        } else {
            dashboard.koDataMonthlyExpensesVal(kendo.toString(Data, "n2"))
        }

        var DataPrev = res.Data[1]
        if (DataPrev < 0) {
            dashboard.koDataMonthlyExpPrev("(" + kendo.toString(DataPrev * -1, "n2") + ")")
        } else {
            dashboard.koDataMonthlyExpPrev(kendo.toString(DataPrev, "n2") + " (Prev Month)")
        }

        callback();
    });
}
dashboard.renderChartMonthlyExpenses = function () {
    data = dashboard.koDataMonthlyRevenue()
    var max = 0
    if (data[0] > data[1]) {
        var value = 10 / 100 * data[0]
        max = data[0] + value
    } else {
        var value = 10 / 100 * data[1]
        max = data[1] + value
    }
    $("#monthlyExpenses").kendoChart({
        seriesColors: ["#f4222d", "#195781"],
        seriesDefaults: {
            type: "bullet"
        },
        legend: {
            visible: false,
        },
        series: [{
            target: {
                border: {
                    width: 3,
                    color: "#195781"
                }
            },
            data: [dashboard.koDataMonthlyRevenue()],
            labels: {
                visible: true,
                background: "transparent",
                template: "#= kendo.toString(value, n2) #"
            }
        }],
        categoryAxis: {
            majorGridLines: {
                visible: false
            },
            majorTicks: {
                visible: false
            }
        },
        valueAxis: {
            plotBands: [{
                from: 0,
                to: data[1],
                color: "#3498db",
                opacity: .6
            }],
            visible: false,
            majorGridLines: {
                visible: false
            },
            minorGridLines: {
                visible: false
            },
            minorTicks: {
                visible: true
            },
            min: 0,
            max: max,
        },
        tooltip: {
            visible: true,
            template: "Previous: #= kendo.toString(value.target, 'n2') # <br /> Current: #= kendo.toString(value.current, 'n2') #"
        },
        chartArea: {
            height: 60
        },
    });
}
dashboard.shortLabels = function (value) {
    if (value == 0) {
        return "0"
    }
    if (value > 0 && value < 1000000) {
        return kendo.toString((value / 1000), "\\") + "k";
    }
    if (value > 0 && value >= 1000000) {
        return kendo.toString((value / 1000000), "\\") + "M";
    }
    // negatif
    if (value < 0 && value > -1000) {
        return "-" + kendo.toString((value / 1000), "\\") + "k";
    }
    if (value < -1000) {
        return "-" + kendo.toString((value / 1000000), "\\") + "M";
    }
}
dashboard.getDataChartTopExpenses = function(callback){
    model.Processing(true)
    var url = "/dashboard/getdatatopfiveexpenses"
    ajaxPost(url, {}, function(res){
        model.Processing(false)
        dashboard.koData5topExpenses(res.Data)
        callback();
    })
}
dashboard.renderChartFiveTopExpenses = function(){
    var colorGreen = ["#9de219","#90cc38","#068c35","#006634","#004d38","#033939"]
    $("#fiveTopExpenses").kendoChart({
        dataSource: {
            data: dashboard.koData5topExpenses(),
        },
        seriesColors : colorGreen,
        title: {
            position: "top",
            // text: dashboard.koDateMonthlyNow()
        },
        chartArea: {
            // background: "",
            // width: 200,
            height: 300
        },
        legend: {
            visible: true,
            position: "bottom",
            labels: {
                font: "10px Arial,Helvetica,sans-serif"
            }
        },
        theme: "bootstrap",
        // seriesDefaults: {
        //     labels: {
        //         visible: true,
        //         background: "transparent",
        //         template: "#= category #: \n #= value#"
        //     }
        // },
        series: [{
            type: "pie",
            field: "Amount",
            categoryField: "_id",
        }],
        tooltip: {
            visible: true,
            template:"#= category#<br /> #= ChangeToRupiah(value) #"
        }
    })
}
dashboard.getDataChartTopRevenue = function(callback){
    model.Processing(true)
    var url = "/dashboard/getdatatopfiverevenue"
    ajaxPost(url, {}, function(res){
        model.Processing(false)
        for(i in res.Data){
            res.Data[i].Amount =res.Data[i].Amount*-1 
        }
        dashboard.koData5topRevenue(res.Data)
        callback();
    })
}
dashboard.renderChartFiveTopRevenue = function(){
    var colorGreen = ["#9de219","#90cc38","#068c35","#006634","#004d38","#033939"]
    $("#fiveTopRevenue").kendoChart({
        dataSource: {
            data: dashboard.koData5topRevenue(),
        },
        seriesColors : colorGreen,
        title: {
            position: "top",
            // text: dashboard.koDateMonthlyNow()
        },
        legend: {
            visible: true,
            position: "bottom",
            labels: {
                font: "10px Arial,Helvetica,sans-serif"
            }
        },
        chartArea: {
            // background: "",
            // width: 200,
            height: 300
        },
        theme: "bootstrap",
        // seriesDefaults: {
        //     labels: {
        //         visible: true,
        //         background: "transparent",
        //         template: "#= category #: \n #= value#"
        //     }
        // },
        series: [{
            type: "pie",
            field: "Amount",
            categoryField: "_id",
        }],
        tooltip: {
            visible: true,
            template:"#= category#<br /> #= ChangeToRupiah(value) #"
        }
    })
}
dashboard.getDataCurrentAsset= function(callback){
    model.Processing(true)
    var url = "/dashboard/getdatacurrentasset"
    ajaxPost(url, {}, function(res){
        model.Processing(false)
        var Data = res.Data
        var debet = _.sumBy(Data, 'SumDebet');
        var credit = _.sumBy(Data, 'SumCredit');
        var saldo = _.sumBy(Data, 'SumSaldo');
        var newData= [
            {_id:"CASH FLOW DEBET", Amount:debet},
            {_id:"CASH FLOW CREDIT", Amount:credit },
            {_id:"CASH FLOW SALDO", Amount:saldo}
        ]
        if (Data.length==0){
            newData = []
        }          
        dashboard.koDataCurrentAssets(newData)
        callback();
    })
}
dashboard.renderChartCurrentAsset = function(){
    var color = ["#ff0909","#005aff","#9de219","#ff5e00","#00eab4","#033939"]
    $("#currentAsset").kendoChart({
        dataSource: {
            data: dashboard.koDataCurrentAssets(),
        },
        seriesColors : color,
        title: {
            position: "top",
        },
        chartArea: {
            // background: "",
            // width: 200,
            height: 300
        },
        legend: {
            visible: true,
            position: "bottom",
            labels: {
                font: "10px Arial,Helvetica,sans-serif"
            }
        },
        theme: "bootstrap",

        series: [{
            type: "pie",
            field: "Amount",
            categoryField: "_id",
        }],
        tooltip: {
            visible: true,
            template:"#= category#<br /> #= ChangeToRupiah(value) #"
        }
    })
}
dashboard.getDataInvoice =  function(callback){
    model.Processing(true)
    var url = "/dashboard/getdatainvoicedashboard"
    ajaxPost(url, {}, function(res){
        var dataOS = _.filter(res.Data, function(o) {
            return o.Status == "Out Standing";
        });
        var total = _.sumBy(dataOS, 'GrandTotalIDR')
        dashboard.koLenOSInvoice(dataOS.length)
        dashboard.koTotalOSInvoice(ChangeToRupiah(total))
        dashboard.koDataINV(res.Data)
        model.Processing(false)
        callback();
    })
}

dashboard.renderGridInvoice = function(){
        $("#gridInvoice").kendoGrid({
            dataSource:{
                data : dashboard.koDataINV(),
                pageSize: 5,
                sort:({
                    field:"DateCreated",
                    dir: "desc"
                })
            },
            // height: 500,
            pageable: {
                refresh: true,
                pageSizes: true,
                buttonCount: 5
            },
            // width: 140,
            sortable: true,
            scrollable: true,
            columns: [{
                field: 'DateCreated',
                title: 'Date',
                width: 100,
                template: function(e){
                    var date = moment(e.DateCreated).format("DD-MMM-YYYY")
                    if(date=="01-Jan-0001"){
                        date = ""
                    }
                    return date
                }
            }, {
                field: 'DocumentNo',
                title: 'Document Number',
                width: 160,
            }, {
                field: 'Status',
                title: 'Status',
                width: 100,
            },{
                field: 'GrandTotalIDR',
                title: 'Total',
                width: 150,
                template: "#=ChangeToRupiah(GrandTotalIDR)#",
                attributes: {
                    style: "text-align:right;"
                },
            },{
                field: 'AlreadyPaid',
                title: 'Payment',
                width: 150,
                template: "#=ChangeToRupiah(AlreadyPaid)#",
                attributes: {
                    style: "text-align:right;"
                },
            }],
            dataBound: function(e){
                var columns = e.sender.columns;
                var Stat = this.wrapper.find(".k-grid-header [data-field=" + "Status" + "]").index();
                dataView = this.dataSource.view();
                for (var i = 0; i < dataView.length; i++) {
                    if (dataView[i].Status == "PAID") {
                        var row1 = e.sender.tbody.find("[data-uid='" + dataView[i].uid + "']");
                        var cell1 = row1.children().eq(Stat);
                        cell1.addClass('greenBackground')
                    }
                }
            }
        })
}
dashboard.getDataPO =  function(callback){
    model.Processing(true)
    var url = "/dashboard/getdatapodashboard"
    ajaxPost(url, {}, function(res){
        var dataOS = _.filter(res.Data, function(o) {
            return o.Status == "Out Standing";
        });
        var total = _.sumBy(dataOS, 'GrandTotalIDR')
        dashboard.koLenOSPO(dataOS.length)
        dashboard.koTotalOSPO(ChangeToRupiah(total))
        dashboard.koDataINV(res.Data)
        model.Processing(false)
        callback();
    })
}
dashboard.renderGridPO = function(){
    $("#gridPO").kendoGrid({
        dataSource:{
            data : dashboard.koDataINV(),
            pageSize: 5,
            sort:({
                field:"DatePosting",
                dir: "desc"
            })
        },
        // height: 500,
        pageable: {
            refresh: true,
            pageSizes: true,
            buttonCount: 5
        },
        // width: 140,
        sortable: true,
        scrollable: true,
        columns: [{
            field: 'DatePosting',
            title: 'Date',
            width: 100,
            template: function(e){
                var date = moment(e.DatePosting).format("DD-MMM-YYYY")
                if(date=="01-Jan-0001"){
                    date = ""
                }
                return date
            }
        }, {
            field: 'DocumentNumber',
            title: 'Document Number',
            width: 160,
        }, {
            field: 'Status',
            title: 'Status',
            width: 100,
        },{
            field: 'GrandTotalIDR',
            title: 'Total',
            width: 150,
            template: "#=ChangeToRupiah(GrandTotalIDR)#",
            attributes: {
                style: "text-align:right;"
            },
        },{
            field: 'AlreadyPaid',
            title: 'Payment',
            width: 150,
            template: "#=ChangeToRupiah(AlreadyPaid)#",
            attributes: {
                style: "text-align:right;"
            },
        }],
        dataBound: function(e){
            var columns = e.sender.columns;
            var Stat = this.wrapper.find(".k-grid-header [data-field=" + "Status" + "]").index();
            dataView = this.dataSource.view();
            for (var i = 0; i < dataView.length; i++) {
                if (dataView[i].Status == "PAID") {
                    var row1 = e.sender.tbody.find("[data-uid='" + dataView[i].uid + "']");
                    var cell1 = row1.children().eq(Stat);
                    cell1.addClass('greenBackground')
                }
            }
        }
    })
}
dashboard.getDataBANK= function(){
    model.Processing(true)
    var url = "/dashboard/getdatabank"
    ajaxPost(url, {}, function(res){
        model.Processing(false)
        var bankUSD = 0.0
        var bankIDR = 0.0
        var cash = 0.0
        var ptc = 0.0
        for(i in res.Data){
            if (res.Data[i]._id == 1120){
                bankUSD = bankUSD +res.Data[i].Amount
            }else if(res.Data[i]._id == 1121){
                bankIDR = bankIDR +res.Data[i].Amount
            }else if(res.Data[i]._id == 1110){
                cash = cash +res.Data[i].Amount
            }else if(res.Data[i]._id == 1111) {
                ptc = ptc +res.Data[i].Amount
            }
        }
        console.log(bankIDR)
        console.log(res.Data)
        //console.log(ptc)
        dashboard.koIDRamount(ChangeToRupiah(bankIDR))
        dashboard.koUSDamount(ChangeToRupiah(bankUSD))
        dashboard.koPTCamount(ChangeToRupiah(ptc))
        dashboard.koPETTYamount(ChangeToRupiah(cash))
    })
}
dashboard.toDashbord1 = function(){
    dashboard.titleDashboard("Dashboard Summary")
    dashboard.visibleDashboard(false)
}
dashboard.init = function () {
    dashboard.getDataChartRevenue(function () {
        dashboard.renderChartRev()
    })
    dashboard.getDataChartExpenses(function () {
        dashboard.renderChartExp()
    })
    dashboard.getDataChartNetProfit(function () {
        dashboard.renderChartNet()
    })
    dashboard.getDataChartRevExNet(function () {
        dashboard.renderChartRevExNet()
    })
    dashboard.getDataChartMonthlyRevenue(function () {
        dashboard.renderChartMonthlyRevenue()
    })
    dashboard.getDataChartMonthlyExpenses(function () {
        dashboard.renderChartMonthlyExpenses()
    })
    
}
dashboard.renderDashvbord2 = function(){
    dashboard.titleDashboard("Dashboard Summary Detail")
    dashboard.visibleDashboard(true)
    dashboard.getDataChartTopExpenses(function () {
        setTimeout(function(){
            dashboard.renderChartFiveTopExpenses()
        },300);
    })
    dashboard.getDataChartTopRevenue(function () {
        setTimeout(function(){
            dashboard.renderChartFiveTopRevenue()
        },300);
    })
    dashboard.getDataCurrentAsset(function () {
        setTimeout(function(){
            dashboard.renderChartCurrentAsset()
        },300);
    })
    dashboard.getDataBANK()
    dashboard.getDataInvoice(function () {
        dashboard.renderGridInvoice()
    })
    dashboard.getDataPO(function () {
        dashboard.renderGridPO()
    })
    setTimeout(function () {
        $(window).resize(function () {
            $("#fiveTopExpenses svg").width(Number($(window).width()));
            $("#fiveTopExpenses").data("kendoChart").refresh();
            $("#fiveTopRevenue svg").width(Number($(window).width()));
            $("#fiveTopRevenue").data("kendoChart").refresh();
            $("#currentAsset svg").width(Number($(window).width()));
            $("#currentAsset").data("kendoChart").refresh();
        });
    }, 1000);
}
$(window).resize(function () {
    $("#coulumn svg").width(Number($(window).width()));
    $("#coulumn").data("kendoChart").refresh();
});
// http://dojo.telerik.com/Opova
$(function () {
    model.Processing(true)
    // dashboard.renderDashvbord2()
    dashboard.init()
    dashboard.getDateNow()
    dashboard.getDateMonthly()
    setTimeout(function () {
        $(window).resize(function () {
            $("#chartRev svg").width(Number($(window).width()));
            $("#chartRev").data("kendoChart").refresh();
            $("#chart1 svg").width(Number($(window).width()));
            $("#chart1").data("kendoChart").refresh();
            $("#chart2 svg").width(Number($(window).width()));
            $("#chart2").data("kendoChart").refresh();
            $("#monthlyRevenue svg").width(Number($(window).width()));
            $("#monthlyRevenue").data("kendoChart").refresh();
            $("#monthlyExpenses svg").width(Number($(window).width()));
            $("#monthlyExpenses").data("kendoChart").refresh();
            // $("#fiveTopExpenses svg").width(Number($(window).width()));
            // $("#fiveTopExpenses").data("kendoChart").refresh();
            // $("#fiveTopRevenue svg").width(Number($(window).width()));
            // $("#fiveTopRevenue").data("kendoChart").refresh();
            // $("#currentAsset svg").width(Number($(window).width()));
            // $("#currentAsset").data("kendoChart").refresh();
        });
    }, 1000);
})