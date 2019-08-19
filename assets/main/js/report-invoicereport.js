var report= {}
report.DatePageBar = ko.observable()
report.DateStart = ko.observable(moment().startOf('month').format('YYYY-MM-DD'))
report.DateEnd = ko.observable(new Date)
report.dataGridreport = ko.observableArray([])
report.dataMasterSales = ko.observableArray([])
report.dataDropDownSales= ko.observableArray([])
report.valueSalesCode = ko.observable()
report.visibleDDsales = ko.observable(false)
report.dataMasterCustomer = ko.observableArray([])
report.dataDropDownCustomerFilter = ko.observableArray([])
report.valueCustomerCode = ko.observable()
report.visibleDDcustomer = ko.observable(false)
report.dataDDReportType = ko.observableArray([
    {
        "Text": "All",
        "Value": "All",
        "TitleReport": "INVOICE REPORT"
    },
    {
        "Text": "Sales",
        "Value": "Sales",
        "TitleReport": "INVOICE BY SALES"
    },
    {
        "Text": "Customer",
        "Value": "Customer",
        "TitleReport": "INVOICE BY CUSTOMER"
    }
])
report.valueDDReportType = ko.observable('All')
report.dataDDReportType2 = ko.observableArray([
    {
        "Text": "Summary",
        "Value": "Summary",
    },
    {
        "Text": "Detail",
        "Value": "Detail",
    }
])
report.valueDDReportType2 = ko.observable('Summary')
report.getDateNow = function() {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    report.DatePageBar(page)
}
report.getDataSales = function() {
    model.Processing(true)
    ajaxPost('/master/getdatasales', {}, function(res) {
        if (res.IsError){
            model.Processing(false)
            return swal('Error!', res.Mesaage, "errror")
        }
        report.dataMasterSales(res.Data)
        var DataSales = res.Data
        for (i in DataSales) {
            DataSales[i].Kode = DataSales[i].SalesID + ""
            DataSales[i].Name = DataSales[i].SalesID + " - " + DataSales[i].SalesName
        }
        report.dataDropDownSales(DataSales)
        model.Processing(false)
    })
}
report.getDataCustomer = function() {
    model.Processing(true)
    ajaxPost('/transaction/getcustomer', {}, function(res) {
        if (res.IsError){
            model.Processing(false)
            return swal('Error!', res.Mesaage, "errror")
        }
        report.dataMasterCustomer(res.Data)
        var DataCustomer = res.Data
        for (i in DataCustomer) {
            DataCustomer[i].Kode = DataCustomer[i].Kode + ""
            DataCustomer[i].Name = DataCustomer[i].Kode + "-" + DataCustomer[i].Name
        }
        report.dataDropDownCustomerFilter(DataCustomer)
        model.Processing(false)
    })
}
report.onChangeSalesCode = function(value){
    if (value== "All"){
        report.visibleDDcustomer(false)
        report.visibleDDsales(false)
    }else if (value== "Sales"){
        report.visibleDDsales(true)
        report.visibleDDcustomer(false)
    }else if (value== "Customer"){
        report.visibleDDsales(false)
        report.visibleDDcustomer(true)
    }
}
report.getDataGrid = function(callback){
    var valuefilter = ""
    if (report.valueDDReportType()=="Sales"){
        valuefilter = report.valueSalesCode()
    }else if (report.valueDDReportType()=="Customer"){
        valuefilter = report.valueCustomerCode()
    }
    var param = {
        DateStart : moment(report.DateStart()).format("YYYY-MM-DD"),
        DateEnd : moment(report.DateEnd()).format("YYYY-MM-DD"),
        ReportType: report.valueDDReportType2(),
        ReportBy : report.valueDDReportType(),
        ValueFilter : valuefilter
    }
    model.Processing(true)
    ajaxPost("/report/getdatainvoicereport", param, function(res){
        if (res.IsError){
            model.Processing(false)
            return swal('Error!', res.Mesaage, "errror")
        }
        report.dataGridreport(res.Data)
        callback();
        model.Processing(false)
    })
}
report.renderGrid = function(){
    $("#gridreport").html("")
    var columns = [
        {
            field:"Date",
            title:"Date",
            template: "#=moment(Date).format('DD-MMM-YYYY')#",
            // width:"100px",
        },{
            field:"Term",
            title:"Term",
            width:"80px",
        },{
            field:"DueDate",
            title:"DueDate",
            template: "#=moment(DueDate).format('DD-MMM-YYYY')#",
        },{
            field:"DocumentNumber",
            title:"Doc. Number"
        },{
            field:"StoreLocationName",
            title:"Location",
            template: function (dt) {
                return dt.StoreLocationName+" ("+dt.StoreLocationId+")"
            }
        },{
            field: "CustomerName",
            title: "Curtomer Name"
        },{
            field: "SalesName",
            title: "Sales Name"
        },{
            field: "Total",
            title: "Total",
            template: '#=ChangeToRupiah(FormatCurrency(Total))#',
            attributes: {
                style: "text-align:right;"
            },
        }
        
    ]
    if (report.valueDDReportType2() == "Detail"){
        columns = [
            {
                field:"Date",
                title:"Date",
                template: "#=moment(Date).format('DD-MMM-YYYY')#",
                // width:"100px",
            },{
                field:"Term",
                title:"Term",
                // width:"80px",
            },{
                field:"DueDate",
                title:"DueDate",
                template: "#=moment(DueDate).format('DD-MMM-YYYY')#",
                // width:"100px",
            },{
                field:"DocumentNumber",
                title:"Doc. Number"
            },{
                field:"StoreLocationName",
                title:"Location",
                template: function (dt) {
                    return dt.StoreLocationName+" ("+dt.StoreLocationId+")"
                }
            },{
                field: "CustomerName",
                title: "Curtomer Name"
            },{
                field: "SalesName",
                title: "Sales Name"
            },{
                field: "Item",
                title: "Item"
            },{
                field: "Price",
                title: "Amount",
                template: '#=ChangeToRupiah(FormatCurrency(Price))#',
                attributes: {
                    style: "text-align:right;"
                },
            },{
                field: "Qty",
                title: "Qty",
                width:"70px",
                attributes: {
                    style: "text-align:center;"
                },
            },{
                field: "Total",
                title: "Total",
                template: '#=ChangeToRupiah(FormatCurrency(Total))#',
                attributes: {
                    style: "text-align:right;"
                },
            }
        ]
    }
    $("#gridreport").kendoGrid({
        dataSource: {
            data: report.dataGridreport(),
            sort: {
                field: 'Date',
                dir: 'asc',
            }
        },
        excelExport: function(e) {
            // var sheet = e.workbook.sheets[0];
            // for (var rowIndex = 1; rowIndex < sheet.rows.length; rowIndex++) {
            //     var row = sheet.rows[rowIndex];
            //     for (var cellIndex = 0; cellIndex < row.cells.length; cellIndex ++) {
            //         if (cellIndex == 0 ||cellIndex == 2 ){
            //             row.cells[cellIndex].format = "yy-MM-dd"
            //         }
            //     }
            // }
            ProActive.kendoExcelRender(e, "report-invoice", function(row, sheet){
                for(var ci = 0; ci < row.cells.length; ci++){
                    var cell = row.cells[ci];
                    if (row.type == "data"){
                        if (ci == 0 || ci == 2 ) {
                            cell.format = "yyyy-MMM-dd";
                            cell.value = moment(cell.value).format("DD-MMM-YYYY")
                        }
                    }
                }
            });
        },
        height: 500,
        width: 140,
        // filterable: true,
        scrollable: true,
        columns: columns
    })
}
report.exportExcel = function(){
    $("#gridreport").getKendoGrid().saveAsExcel();
}
report.ExportToPdf = function(){
    var valuefilter = ""
    if (report.valueDDReportType()=="Sales"){
        valuefilter = report.valueSalesCode()
    }else if (report.valueDDReportType()=="Customer"){
        valuefilter = report.valueCustomerCode()
    }
    var param = {
        DateStart : moment(report.DateStart()).format("YYYY-MM-DD"),
        DateEnd : moment(report.DateEnd()).format("YYYY-MM-DD"),
        ReportType: report.valueDDReportType2(),
        ReportBy : report.valueDDReportType(),
        ValueFilter : valuefilter
    }
    model.Processing(true)
    ajaxPost("/report/exporttopdfreportinvoice", param, function(res){
        if (res.IsError){
            model.Processing(false)
            return swal('Error!', res.Mesaage, "errror")
        }
        window.open('/res/docs/report/pdf/' + res.Data, '_blank');
        model.Processing(false)
    })
}
report.detailReportPdf = function(){
    var valuefilter = ""
    if (report.valueDDReportType()=="Sales"){
        valuefilter = report.valueSalesCode()
    }else if (report.valueDDReportType()=="Customer"){
        valuefilter = report.valueCustomerCode()
    }
    var param = {
        DateStart : moment(report.DateStart()).format("YYYY-MM-DD"),
        DateEnd : moment(report.DateEnd()).format("YYYY-MM-DD"),
        ReportType: report.valueDDReportType2(),
        ReportBy : report.valueDDReportType(),
        ValueFilter : valuefilter
    }
    model.Processing(true)
    ajaxPost("/report/detailreportpdf", param, function(res){
        if (res.IsError){
            model.Processing(false)
            return swal('Error!', res.Mesaage, "errror")
        }
        window.open('/res/docs/report/pdf/' + res.Data, '_blank');
        model.Processing(false)
    })
}
report.refreshGrid = function(){
    report.getDataGrid(function(){
        report.renderGrid()
    })
}

// report.initDate = function () {
//     var dtpStart = $('#dateStart').data('kendoDatePicker');
//     var dtpEnd = $('#dateEnd').data('kendoDatePicker');
//     dtpStart.value(moment().startOf('month').toDate());
//     dtpEnd.value(moment().startOf('day').toDate());
//     dtpStart.max(dtpEnd.value());
//     dtpEnd.min(dtpStart.value());

//     dtpStart.bind("change", function () {
//         dtpEnd.min(dtpStart.value());
//     });
//     dtpEnd.bind("change", function () {
//         dtpStart.max(dtpEnd.value());
//     });
// }
report.onChangeDateStart = function(val){
    if (val.getTime()>report.DateEnd().getTime()){
        report.DateEnd(val)
    }
}

report.filterText = function(term) {
    var filter = { logic: "or", filters: [] };
    var filteredFields = ["DocumentNumber", "SalesName", "CustomerName", "StoreLocationName"]
    $searchValue = term || $("#textSearch").val();
    if ($searchValue) {
        for (var k in filteredFields)
            filter.filters.push({ field: filteredFields[k], operator:"contains", value:$searchValue});
    } 
    $("#gridreport").data("kendoGrid").dataSource.query({ filter: filter });
}

report.init = function(){
    report.getDateNow()
    // report.initDate()
    report.getDataSales()
    report.getDataCustomer()
    report.getDataGrid(function(){
        report.renderGrid()
    })
}
$(function() {
    report.init()
    $("#textSearch").on("keyup blur change", function () {
        report.filterText();
    });
})