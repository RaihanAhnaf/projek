model.Processing(false)
model.Processing(false)
var opex = {}
opex.DateStart = ko.observable(moment().startOf('month').format('YYYY-MM-DD hh:mm'))
opex.DateEnd = ko.observable(new Date)
opex.koFilterIsActive = ko.observable(false)
opex.koDataOpex = ko.observableArray([])
opex.TitelFilter = ko.observable(" Hide Filter")
opex.DatePageBar = ko.observable()
opex.typeFilter = ko.observable("")
opex.textSearch = ko.observable()
opex.filterindicator = ko.observable(false)
opex.valueDepartment = ko.observable([])
opex.valueSales = ko.observable([])
opex.titleOpex= ko.observable("Operating Expenses Summary")
opex.visibleSummary = ko.observable(true)
opex.dataTop10Opex = ko.observableArray([])
opex.dataTop10OpexDetail = ko.observableArray([])
opex.DataOpexOneYears = ko.observableArray([])
opex.categoryData= ko.observableArray([])
opex.filterByValue = ko.observable('Detail')
opex.dataOpexPerPeriode = ko.observableArray([])
opex.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    opex.DatePageBar(page)
}

opex.refreshGrid = function () {
    if(opex.valueDepartment()!= ""){
        $("#departmenDropdown-list > span > input").val("")
    }

    if(opex.valueSales()!= ""){
        $("#salesDropdown-list > span > input").val("")
    }

    opex.koFilterIsActive(true)
    if (opex.filterByValue()== "Detail"){
        setTimeout(function(){
            opex.getDataOpex(function () {
                opex.renderGrid()
            })
        },200)
    }else{
        // $('#gridOpex').html("")
        setTimeout(function(){
            opex.getDataOpexPeriode(function () {
                opex.renderGridOpexPeriode()
            })
        },200)
    } 
    opex.getDataTopTenOpex(function(){
        opex.renderPieChartOpex()
    })
    opex.getDataDetailTopTenOpex(function(){
        opex.renderGridDetailTopTenOpex()
    })
    opex.getDataChartColumn(function(){
        opex.renderChartRevExNet()
    })
}
opex.getDataOpex = function (callback) {
    model.Processing(true)
    var url = "/report/getdataopex"
    var dateStart = moment($('#dateStart').data('kendoDatePicker').value()).format("YYYY-MM-DD")
    var dateEnd = moment($('#dateEnd').data('kendoDatePicker').value()).format("YYYY-MM-DD")
    var DepartmentContainstr = $("#departmenDropdown-list > span > input").val()
    if (DepartmentContainstr == undefined){
        DepartmentContainstr = ""
    }
    var SalesContainstr = $("#salesDropdown-list > span > input").val()
    if (SalesContainstr == undefined){
        SalesContainstr = ""
    }
    var param = {}
    if (opex.koFilterIsActive() == true) {
         param = {
             DateStart: dateStart,
             DateEnd: dateEnd,
             Filter: true,
             TextSearch: opex.textSearch(),
             Department : opex.valueDepartment(),
             DepartmentContains :DepartmentContainstr.toUpperCase(),
             SalesCode : opex.valueSales(),
             SalesContains :SalesContainstr.toUpperCase()
         }
     } else {
         param = {
            DateStart: dateStart,
            DateEnd: dateEnd,
            Filter: false,
            Department : opex.valueDepartment(),
            DepartmentContains :DepartmentContainstr.toUpperCase(),
            SalesCode : opex.valueSales(),
            SalesContains :SalesContainstr.toUpperCase()
         }
     }
    ajaxPost(url, param, function (res) {
        if (res.IsError) {
            swal("Search Not Found!", res.Message, "warning")
            $('#textSearch').val("")
            return
        }
        var data = res.Data
        var totalAmount = _.sumBy(data, function (e) {
            return e.Amount
        })
        for (i in data) {
            if (totalAmount == 0) {
                var tempPercentage = 0
            } else {
                var tempPercentage = (data[i].Amount / totalAmount) / 100 * 100
            }
            data[i].Percentage = tempPercentage
        }
        opex.koDataOpex(data)
        model.Processing(false)
        callback()
    })
}
opex.renderGrid = function () {
    $('#gridOpex').html("")
    var data = opex.koDataOpex()
    var columns = [{
        field: 'Acc_Code',
        title: 'Acc Code',
        width: 50
    }, {
        field: 'Acc_Name',
        title: 'Account Name',
        width: 200,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>Total Operating Expenses:</div>",
    }, {
        field: 'Amount',
        title: 'Amount',
        headerAttributes: {
            style: "text-align:right"
        },
        width: 50,
        template: "#=ChangeToRupiah(Amount)#",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }, {
        field: 'Percentage',
        title: '%',
        headerAttributes: {
            style: "text-align:right"
        },
        width: 20,
        template: "#=kendo.toString(Percentage, 'p1')#",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= kendo.toString(sum, 'p0') #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }]
    $('#gridOpex').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
                field: "Amount",
                aggregate: "sum"
            }, {
                field: "Percentage",
                aggregate: "sum"
            }],
            sort: {
                field: 'Acc_Code',
                dir: 'asc',
            }
        },
        excel: {
            fileName: "report-opex.xlsx"
        },
        excelExport: function (op) {
            var rows = op.workbook.sheets[0].rows;
            for (var ri = 0; ri < rows.length; ri++) {
                var row = rows[ri];
                for (var ci = 0; ci < row.cells.length; ci++) {
                    var cell = row.cells[ci];
                    if (row.type == "group-footer" || row.type == "footer") {
                        var text = $(cell.value).text();
                        cell.hAlign = "right";
                        cell.value = text
                    }
                    if (ci == 3) {
                        if (row.type == "data") {
                            cell.value = kendo.toString(cell.value, 'p0')
                            cell.hAlign = "right";
                        }
                    }
                    if (ci == 2) {
                        if (row.type == "data") {
                            cell.value = parseFloat(cell.value)
                            cell.format = "#,##0.00_);(#,##0.00);0.00;"
                            // Set the alignment
                            cell.hAlign = "right";
                        }
                    }
                }
            }
        },
        height: 550,
        sortable: true,
        scrollable: true,
        columns: columns,
        detailInit: detailInit,
    })
}

function detailInit(e) {
    var DepartmentContainstr = $("#departmenDropdown-list > span > input").val()
    if (DepartmentContainstr == undefined){
        DepartmentContainstr = ""
    }
    var SalesContainstr = $("#salesDropdown-list > span > input").val()
    if (SalesContainstr == undefined){
        SalesContainstr = ""
    }
    var url = "/report/getdatadetailopex"
    var dateStart = moment($('#dateStart').data('kendoDatePicker').value()).format("YYYY-MM-DD")
    var dateEnd = moment($('#dateEnd').data('kendoDatePicker').value()).format("YYYY-MM-DD")
    var param = {
        DateStart: dateStart,
        DateEnd: dateEnd,
        Filter: opex.koFilterIsActive(),
        AccCode: e.data.Acc_Code,
        Department : opex.valueDepartment(),
        DepartmentContains :DepartmentContainstr.toUpperCase(),
        SalesCode : opex.valueSales(),
        SalesContains :SalesContainstr.toUpperCase()
    }
    if (e.data.Amount != 0) {
        model.Processing(true)
        ajaxPost(url, param, function (res) {
            model.Processing(false)
            $("<div/>").appendTo(e.detailCell).kendoGrid({
                columns: [{
                    title: "Date",
                    field: "PostingDate",
                    template: function (e) {
                        var value = e.PostingDate
                        return moment(value).format("DD MMM YYYY")
                    },
                    width: 120
                }, {
                    title: "Document Number",
                    field: "DocumentNumber",
                    width: 170
                }, {
                    title: "Account Code",
                    field: "Acc_Code",
                    width: 100
                }, {
                    title: "Account Name",
                    field: "Acc_Name",
                    width: 150
                }, {
                    title: "Department",
                    field: "Department",
                    width: 100
                }, {
                    title: "Sales",
                    field: "SalesName",
                    width: 100
                }, {
                    title: "Description",
                    field: "Description",
                    width: 150
                }, {
                    title: "Amount",
                    field: "Amount",
                    template: "#=ChangeToRupiah(Amount)#",
                    width: 100,
                    attributes: {
                        style: "text-align:right;"
                    }
                }, ],
                dataSource: res.Data
            });
        });
    }

    // http://dojo.telerik.com/aliJoG
}
opex.getDataOpexPeriode = function(callback){
    model.Processing(true)
    var url = "/report/getdataopexperiode"
    var dateStart = moment($('#dateStart').data('kendoDatePicker').value()).format("YYYY-MM-DD")
    var dateEnd = moment($('#dateEnd').data('kendoDatePicker').value()).format("YYYY-MM-DD")
    var DepartmentContainstr = $("#departmenDropdown-list > span > input").val()
    if (DepartmentContainstr == undefined){
        DepartmentContainstr = ""
    }
    var SalesContainstr = $("#salesDropdown-list > span > input").val()
    if (SalesContainstr == undefined){
        SalesContainstr = ""
    }
    var param = {}
    if (opex.koFilterIsActive() == true) {
         param = {
             DateStart: dateStart,
             DateEnd: dateEnd,
             Filter: true,
             TextSearch: opex.textSearch(),
             Department : opex.valueDepartment(),
             DepartmentContains :DepartmentContainstr.toUpperCase(),
             SalesCode : opex.valueSales(),
             SalesContains :SalesContainstr.toUpperCase()
         }
     } else {
         param = {
            DateStart: dateStart,
            DateEnd: dateEnd,
            Filter: false,
            Department : opex.valueDepartment(),
            DepartmentContains :DepartmentContainstr.toUpperCase(),
            SalesCode : opex.valueSales(),
            SalesContains :SalesContainstr.toUpperCase()
         }
     }
    ajaxPost(url, param, function (res) {
        if (res.IsError) {
            swal("Search Not Found!", res.Message, "warning")
            $('#textSearch').val("")
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
            var m = v._id.Month
            var y = v._id.Year 
            var dateMy = new Date(y,m-1, 01)
            var my = moment(dateMy).format("MMMYYYY")
            DataAmount.push({
                Acc_Code :v._id.Acc_Code,
                Acc_Name :v._id.Acc_Name,
                MonthYear : my,
                Amount : v.Amount
            })
        })
        // console.log(DataAmount)
        var dataAmont2 =_.chain(DataAmount).groupBy("Acc_Code").toPairs().map(function (currentItem) {
            return _.zipObject(["Acc_Code", "DataItem"], currentItem);
        }).value();
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
                    Acc_Name:v.Acc_Name,
                    DataItem : itemMy,
                    TotalAmount: 0
                })
            }
        })
        opex.dataOpexPerPeriode(dataAcc)
        model.Processing(false)
        callback()
    })
}
opex.renderGridOpexPeriode = function(){
    $('#gridOpex').html("")
    var data = opex.dataOpexPerPeriode()
    var columns = [{
        title: 'Account Name',
        field: 'Acc_Name',
        width: 200,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>Total Operating Expenses:</div>",
    }]
    var dateStart = moment($('#dateStart').data('kendoDatePicker').value()).format("YYYY-MM-DD")
    var dateEnd = moment($('#dateEnd').data('kendoDatePicker').value()).format("YYYY-MM-DD")
    var Aggregate = []
    var idxDate = []
    for (var d = new Date(dateStart); d <= new Date(dateEnd); d.setMonth(d.getMonth() + 1)) {
        var namefield = "DataItem." +moment(d).format("MMMYYYY")
        var date1 = moment(d).format("MMMYYYY")
        idxDate.push(date1)
        //console.log(namefield)
        columns.push({
            title: moment(d).format("MMM YYYY"),
            field: namefield,
            template: "#=ChangeToRupiah("+namefield+")#",
            width: 100,
            attributes: {
                style: "text-align:right;"
            },
            footerTemplate: "<div class='align-right'>#=ChangeToRupiah(sum)#</div>",
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
    // console.log(columns)
    $('#gridOpex').kendoGrid({
        dataSource: {
            data: data,
            aggregate : Aggregate,
            sort: {
                field: 'Acc_Code',
                dir: 'asc',
            }
        },
        excel: {
            fileName: "report-opex-periode.xlsx"
        },
        excelExport: function (op) {
            // op.sender.showColumn("TotalAmount");
            // // op.preventDefault();
            // // setTimeout(function () {
            // //     op.sender.saveAsExcel();
            // // });
            // // // for(i in op.sender.columns){
            // // //     if(op.sender.columns[i].hidden){
            // // //         op.sender.showColumn(i);
            // // //     }
            // // // }
            
            // // // var exportFlag = false;
            // // // if (!exportFlag) {
            // //     op.sender.hideColumn(2);
            // //     op.preventDefault();
            // //     // exportFlag = true;
            // // //     setTimeout(function () {
            // // //         op.sender.saveAsExcel();
            // // //     });
            // // // } else {
            // // //     op.sender.showColumn(1);
            // // //     exportFlag = false;
            // // // }
            var rows = op.workbook.sheets[0].rows;
            for (var ri = 0; ri < rows.length; ri++) {
                var row = rows[ri];
                for (var ci = 0; ci < row.cells.length; ci++) {
                    var cell = row.cells[ci];
                    if (row.type == "group-footer" || row.type == "footer") {
                        var text = $(cell.value).text();
                        cell.hAlign = "right";
                        cell.value = text
                    }
                    if(row.type=="header"){
                        cell.background="#f3de26",
                        cell.color="#000000"
                    }
                    // if (ci == 3) {
                    //     if (row.type == "data") {
                    //         cell.value = kendo.toString(cell.value, 'p0')
                    //         cell.hAlign = "right";
                    //     }
                    // }
                    if (ci >= 1) {
                        if (row.type == "data") {
                            cell.value = parseFloat(cell.value)
                            cell.format = "#,##0.00_);(#,##0.00);0.00;"
                            // Set the alignment
                            cell.hAlign = "right";
                        }
                    }
                }
            }
        },
        height: 550,
        sortable: true,
        scrollable: true,
        columns: columns,
    })
}
opex.setDate = function () {
    var datepicker = $("#dateStart").data("kendoDatePicker");
    var now = new Date();
    var years = moment(now).format("YYYY")
    var Month = moment(now).format("MM")
    var newDate = years + "-" + Month + "-" + "01"
    datepicker.value(new Date(newDate))
}
opex.departmentDropdown = function(){
    $("#departmenDropdown").html("")
    var data = []
    ajaxPost("/transaction/getdatadepartment", {}, function(res){
        var Data = res.Data
        Data.push({
            DepartmentCode:"DEPT/ALL/0001",
            DepartmentName:"ALL",
            _id:"all"
        })
        Data =_.sortBy(Data, function(e){return e.DepartmentName})
        $("#departmenDropdown").kendoMultiSelect({
            filter: "contains",
            dataTextField: "DepartmentName",
            dataValueField: "DepartmentName",
            dataSource: Data,
            placeholder:'Select one',
            change:function(e){
                var multiselect = $("#departmenDropdown").data("kendoMultiSelect");
                // console.log(multiselect.value())
                opex.valueDepartment(multiselect.value())
            },
        });
    })
}
opex.selectedDepartment = function(){
    var dropdownlist = $("#departmenDropdown").data("kendoMultiSelect");
    dropdownlist.value([]);
    dropdownlist.trigger("change");
    opex.selectedSales()
}

opex.salesDropdown = function(){
    $("#salesDropdown").html("")
    var data = []
    ajaxPost("/master/getdatasales", {}, function(res){
        var Data = res.Data
        Data =_.sortBy(Data, function(e){return e.SalesName})
        Data.reverse()
        for (i in Data) {
            Data[i].SalesCode = Data[i].SalesID
            Data[i].SalesName = Data[i].SalesID + " - "+ Data[i].SalesName   
        }
        Data.push({
            SalesCode:"ALL",
            SalesName:"ALL",
            _id:"all"
        })
        Data.reverse()
        $("#salesDropdown").kendoMultiSelect({
            filter: "contains",
            dataTextField: "SalesName",
            dataValueField: "SalesCode",
            dataSource: Data,
            placeholder:'Select one',
            change:function(e){
                var multiselect = $("#salesDropdown").data("kendoMultiSelect");
                // console.log(multiselect.value())
                opex.valueSales(multiselect.value())
            },
        });
    })
}
opex.selectedSales = function(){
    var dropdownlist = $("#salesDropdown").data("kendoMultiSelect");
    dropdownlist.value([]);
    dropdownlist.trigger("change");
}

opex.dowloadExcel = function () {
    // if(opex.valueDepartment()==""){
    //     $("#gridOpex").getKendoGrid().saveAsExcel();
    // }else{
    if (opex.filterByValue()== "Detail"){
        var DepartmentContainstr = $("#departmenDropdown-list > span > input").val()
        if (DepartmentContainstr == undefined){
            DepartmentContainstr = ""
        }
        var SalesContainstr = $("#salesDropdown-list > span > input").val()
        if (SalesContainstr == undefined){
            SalesContainstr = ""
        }
        var url = "/report/exportdetailopextoexcell"
        var dateStart = moment($('#dateStart').data('kendoDatePicker').value()).format("YYYY-MM-DD")
        var dateEnd = moment($('#dateEnd').data('kendoDatePicker').value()).format("YYYY-MM-DD")
        var param = {
            DateStart: dateStart,
            DateEnd: dateEnd,
            Department : opex.valueDepartment(),
            DepartmentContains :DepartmentContainstr.toUpperCase(),
            SalesCode : opex.valueSales(),
            SalesContains :SalesContainstr.toUpperCase()
        }
        model.Processing(true)
        ajaxPost(url, param, function(e){
            model.Processing(false)
            var pom = document.createElement('a');
            pom.setAttribute('href', "/res/docs/report/excel/" + e.Data);
            pom.setAttribute('download', e.Data);
            pom.click();
        })
    }else{
        var exportFlag = true;
        $("#gridOpex").data("kendoGrid").bind("excelExport", function (e) {
            if (exportFlag) {
            e.sender.showColumn("TotalAmount");
            e.preventDefault();
            exportFlag = false;
            e.sender.saveAsExcel();
            } else {
            e.sender.hideColumn("TotalAmount");
            exportFlag = true;
            }
        });
        $("#gridOpex").getKendoGrid().saveAsExcel();
        opex.renderGridOpexPeriode()
    }
    // }
}
opex.ExportToPdf = function(){
    var DepartmentContainstr = $("#departmenDropdown-list > span > input").val()
    if (DepartmentContainstr == undefined){
        DepartmentContainstr = ""
    }
    var SalesContainstr = $("#salesDropdown-list > span > input").val()
    if (SalesContainstr == undefined){
        SalesContainstr = ""
    }
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var url = "/report/exportpdfopex"
    var param = {
        Department : opex.valueDepartment(),
        DateStart: dateStart,
        DateEnd: dateEnd,
        DepartmentContains :DepartmentContainstr.toUpperCase(),
        SalesCode : opex.valueSales(),
        SalesContains :SalesContainstr.toUpperCase()
    }
    model.Processing(true)
    ajaxPost(url, param, function(e){
        window.open('/res/docs/report/pdf/' + e.Data, '_blank');
        model.Processing(false)
    })
}
opex.goToChart = function(){
    opex.visibleSummary(true)
    opex.titleOpex("Operating Expenses Summary")
}
opex.goToSummary = function(){
    opex.visibleSummary(false)
    opex.titleOpex("Operating Expenses")
    // opex.refreshGrid()
}
opex.getDataTopTenOpex = function(callback){
    var DepartmentContainstr = $("#departmenDropdown-list > span > input").val()
    if (DepartmentContainstr == undefined){
        DepartmentContainstr = ""
    }
    var SalesContainstr = $("#salesDropdown-list > span > input").val()
    if (SalesContainstr == undefined){
        SalesContainstr = ""
    }
    var dateStart = moment($('#dateStart').data('kendoDatePicker').value()).format("YYYY-MM-DD")
    var dateEnd = moment($('#dateEnd').data('kendoDatePicker').value()).format("YYYY-MM-DD")
    var url = "/report/getdatatoptenopexpie"
    var param = {
        DepartmentContains : DepartmentContainstr,
        Department : opex.valueDepartment(),
        SalesCode : opex.valueSales(),
        SalesContains :SalesContainstr.toUpperCase(),
        DateStart: dateStart,
        DateEnd: dateEnd,
    }
    model.Processing(true)
    ajaxPost(url, param, function(e){
        model.Processing(false)
        opex.dataTop10Opex(e.Data)
        callback();
    })
}
opex.renderPieChartOpex = function(){
    var colorGreen = ["#9de219","#90cc38","#699627","#50b11e","#068c35","#006634","#004d38","#033939","#0a3507","#1cb510","#054800"]
    $("#TopTenOpexPie").kendoChart({
        dataSource: {
            data: opex.dataTop10Opex(),
        },
        seriesColors : colorGreen,
        title: {
            position: "top",
            // text: dashboard.koDateMonthlyNow()
        },
        chartArea: {
            background: "",
            // width: 200,
            height: 400
        },
        legend: {
            visible: true,
            position: "left",
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
            categoryField: "Account",
        }],
        tooltip: {
            visible: true,
            template:"#= category#<br /> #= ChangeToRupiah(value) #"
        }
    })
}
opex.getDataDetailTopTenOpex = function(callback){
    var DepartmentContainstr = $("#departmenDropdown-list > span > input").val()
    if (DepartmentContainstr == undefined){
        DepartmentContainstr = ""
    }
    var SalesContainstr = $("#salesDropdown-list > span > input").val()
    if (SalesContainstr == undefined){
        SalesContainstr = ""
    }
    var dateStart = moment($('#dateStart').data('kendoDatePicker').value()).format("YYYY-MM-DD")
    var dateEnd = moment($('#dateEnd').data('kendoDatePicker').value()).format("YYYY-MM-DD")
    var url = "/report/getdatatoptenopexdetailgrid"
    var param = {
        DepartmentContains : DepartmentContainstr,
        Department : opex.valueDepartment(),
        SalesCode : opex.valueSales(),
        SalesContains :SalesContainstr.toUpperCase(),
        DateStart: dateStart,
        DateEnd: dateEnd,
    }
    model.Processing(true)
    ajaxPost(url, param, function(e){
        model.Processing(false)
        opex.dataTop10OpexDetail(e.Data)
        callback();
    })
}
opex.renderGridDetailTopTenOpex = function(){
     $('#TopTenOpexGrid').html("")
    var data = opex.dataTop10OpexDetail()
    var columns = [{
        title: "Date",
        field: "DateStr",
        width: 100
    }, {
        field: 'Acc_Name',
        title: 'Account Name',
        width: 100,
    }, {
        title: "Department",
        field: "Department",
        width: 100
    },{
        title: "Sales",
        field: "SalesName",
        width: 100
    }, {
        title: "Description",
        field: "Description",
        width: 150,
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>Total Operating Expenses:</div>",
    },{
        field: 'Amount',
        title: 'Amount',
        headerAttributes: {
            style: "text-align:right"
        },
        width: 100,
        template: "#=ChangeToRupiah(Amount)#",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
        attributes: {
            style: "text-align:right;"
        }
    }]
    $('#TopTenOpexGrid').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
                field: "Amount",
                aggregate: "sum"
            }, {
                field: "Percentage",
                aggregate: "sum"
            }],
            sort: {
                field: 'Acc_Code',
                dir: 'asc',
            }
        },
        excel: {
            fileName: "report-opex.xlsx"
        },
        height: 400,
        sortable: true,
        scrollable: true,
        columns: columns,
    })
}
opex.getDataChartColumn = function (callback) {
    var url = "/report/getdatachartopex"
    var DepartmentContainstr = $("#departmenDropdown-list > span > input").val()
    if (DepartmentContainstr == undefined){
        DepartmentContainstr = ""
    }
    var SalesContainstr = $("#salesDropdown-list > span > input").val()
    if (SalesContainstr == undefined){
        SalesContainstr = ""
    }
    var param = {
        Year: parseInt(moment().format("YYYY")),
        Department : opex.valueDepartment(),
        DepartmentContains :DepartmentContainstr.toUpperCase(),
        SalesCode : opex.valueSales(),
        SalesContains :SalesContainstr.toUpperCase()
    }
    model.Processing(true)
    ajaxPost(url, param, function (res) {
        var Data = res.Data
        Data = _.sortBy(Data, [function (o) {
            return o._id;
        }]);
        for (i in Data) {
            Data[i].Date = new Date(moment().format("YYYY") + "-" + Data[i].Month + "-" + "1"),
            Data[i].Category =  moment(Data[i].Date).format("MMM-YYYY")
        }
        var Ctg = _.sortBy(Data, [function (o) {
            return o.Date;
        }]);
        var Uniq =_.uniqBy(Ctg, function(e){return e.Type}); 
        //console.log(Uniq)
        var CategotyData = _.map(Ctg, 'Category');
        CategotyData =_.uniq(CategotyData)
        var NameSeries = _.map(Uniq, 'Type');
        // console.log(NameSeries)
        var series = []

        for (i in NameSeries){
            var name = NameSeries[i]
            var data = []
            var filtered = _.filter(Ctg, function(o) { return o.Type== name});
            filtered = _.sortBy(filtered, [function (o) {return o.Date;}]);
            for (j in filtered){
                data.push(filtered[j].Amount)
            }
            series.push({
                data : data,
                name : name,
                type: "column"
            })
        }
        // var CategotyData = _.groupBy(Ctg, function(e){
        //     return e.Category
        // })
        // console.log(CategotyData,series)
        opex.categoryData(CategotyData)
        opex.DataOpexOneYears(series)

        callback()
        model.Processing(false)
    })
}
opex.renderChartRevExNet = function () {
    var blueColor = ["#8bc6ed", "#5497c4", "#2e6c96", "#1b4f72", "#1072b5", "#0b8ce5", "#31c5f7", "#042c44", "#023d60", "#00507f"]
    var colorGreen = ["#9de219","#90cc38","#699627","#50b11e","#068c35","#006634","#004d38","#033939","#1cb510","#820000","#ffee51",
    "#8bc6ed", "#5497c4", "#2e6c96", "#1b4f72", "#1072b5", "#0b8ce5", "#ac6800", "#fd4326", "#ff0808", "#4656d4",
    "#508991", "#A39BA8", "#172A3A","#004346","#09BC8A"]
    var colorTes =["#508991", "#A39BA8", "#172A3A","#004346","#09BC8A"]
    var Data = new kendo.data.DataSource({
        data: opex.DataOpexOneYears(),
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

    $("#TopTenOpexColumn").kendoChart({
        // dataSource: Data,
        legend: {
            position: "bottom",

        },
        // seriesDefaults: {
        //     type: "column"
        // },
        title: {
            // text: dashboard.koDateNow(),
            font: "bold 15px Arial,Helvetica,sans-serif",
            color: "black"
        },
        seriesColors: colorGreen,
        series : opex.DataOpexOneYears(),
        // series: [{
        //     name: "#= group.value #",
        //     field: "Amount",
        //     overlay: {
        //         gradient: "none"
        //     },
        // }],
        valueAxis: {
            majorGridLines: {
                visible: false
            },
            minorGridLines: {
                visible: false
            },
            labels: {
                template: "#= opex.shortLabels(value) #"
            }
        },
        theme: "bootstrap",
        categoryAxis: {
            categories: opex.categoryData(),
            majorGridLines: {
                    visible: false
            },
            labels: {
                font:"11px sans-serif",
                // rotation: {
                //     angle: 45,
                //     align: "center"
                // },
            }
        },
        // categoryAxis: {
        //     majorGridLines: {
        //         visible: false
        //     },
        //     minorGridLines: {
        //         visible: false
        //     },
        //     field: "Month",
        //     labels: {
        //         font: "11px sans-serif",
        //         // format: "MMM\nyyyy"
        //     }
        // },
        tooltip: {
            visible: true,
            template: "#= series.name# : #= ChangeToRupiah(value)#",
            theme: "bootstrap",
        },
        chartArea: {
            height: 353,
        },
    });
}
opex.shortLabels = function (value) {
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
opex.search = function(e) {
 opex.textSearch(e)
 opex.filterindicator(true)
 opex.getDataOpex (function() {
     opex.renderGrid()
 })
}
opex.renderFilterBy= function(){
    var Data = [
    {"value":"Detail", "text":"Detail"},
    {"value":"Periode", "text":"Period"}
    ]
    $("#filterBy").kendoDropDownList({
        filter: "contains",
        dataTextField: "text",
        dataValueField: "value",
        dataSource: Data,
        index: 1,
        optionLabel:'Select one',
        change:function(e){
            var dataitem = this.dataItem();
           opex.filterByValue(dataitem.value)
        },
    });
}

opex.onChangeDateStart = function(val){
    if (val.getTime()>opex.DateEnd().getTime()){
        opex.DateEnd(val)
    }
}

opex.init = function () {
    // opex.setDate()
    opex.getDateNow()
    opex.departmentDropdown()
    opex.salesDropdown()
    opex.renderFilterBy()
    opex.getDataTopTenOpex(function(){
        opex.renderPieChartOpex()
    })
    opex.getDataDetailTopTenOpex(function(){
        opex.renderGridDetailTopTenOpex()
    })
    opex.getDataChartColumn(function(){
        opex.renderChartRevExNet()
    })
    setTimeout(function(){
        opex.getDataOpex(function () {
        opex.renderGrid()
    })
    },200)
}
$(document).ready(function () {
    opex.init()
    setTimeout(function () {
        $(window).resize(function () {
            $("#TopTenOpexPie svg").width(Number($(window).width()));
            $("#TopTenOpexPie").data("kendoChart").refresh();
            $("#TopTenOpexColumn svg").width(Number($(window).width()));
            $("#TopTenOpexColumn").data("kendoChart").refresh();
        });
    }, 1000);
})