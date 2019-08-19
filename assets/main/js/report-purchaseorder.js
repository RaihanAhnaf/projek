var report= {}
report.DatePageBar = ko.observable()
report.DateStart = ko.observable(new Date)
report.DateEnd = ko.observable(new Date)
report.dataGridreport = ko.observableArray([])
report.dataMasterSales = ko.observableArray([])
report.dataDropDownSales= ko.observableArray([])
report.valueSalesCode = ko.observable()
report.visibleDDsales = ko.observable(false)
report.dataMasterSupplier = ko.observableArray([])
report.dataDropDownSupplierFilter = ko.observableArray([])
report.valueSupplierCode = ko.observable()
report.visibleDDcustomer = ko.observable(false)
report.isInventory = ko.observable(false)
report.dataDDReportType = ko.observableArray([
    {
        "Text": "All",
        "Value": "All",
        "TitleReport": "PURCHASE ORDER REPORT"
    },
    {
        "Text": "Supplier",
        "Value": "Supplier",
        "TitleReport": "PURCHASE ORDER BY SUPPLIER"
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
report.getDataSupplier = function() {
    model.Processing(true)
    ajaxPost('/master/getdatasupplier', {}, function(res) {
        if (res.IsError){
            model.Processing(false)
            return swal('Error!', res.Mesaage, "errror")
        }
        report.dataMasterSupplier(res.Data)
        var DataCustomer = res.Data
        for (i in DataCustomer) {
            DataCustomer[i].Kode = DataCustomer[i].Kode + ""
            DataCustomer[i].Name = DataCustomer[i].Kode + "-" + DataCustomer[i].Name
        }
        report.dataDropDownSupplierFilter(DataCustomer)
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
    }else if (value== "Supplier"){
        report.visibleDDsales(false)
        report.visibleDDcustomer(true)
    }
}
report.getDataGrid = function(callback){
    var valuefilter = ""
    if (report.valueDDReportType()=="Sales"){
        valuefilter = report.valueSalesCode()
    }else if (report.valueDDReportType()=="Supplier"){
        valuefilter = report.valueSupplierCode()
    }
    var param = {                                                                                                                                                                                                                                                                                                                                                                                                      
        DateStart : moment(report.DateStart()).format("YYYY-MM-DD"),
        DateEnd : moment(report.DateEnd()).format("YYYY-MM-DD"),
        ReportType: report.valueDDReportType2(),
        ReportBy : report.valueDDReportType(),
        ValueFilter : valuefilter,
        IsInventory : report.isInventory()
    }
    model.Processing(true)
    ajaxPost("/report/getdatapurchaseorderreport", param, function(res){
        if (res.IsError){
            model.Processing(false)
            return swal('Error!', res.Mesaage, "error")
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
            width:"100px",
        },{
            field:"DocumentNumber",
            title:"Doc. Number"
        },{
            field: "SupplierName",
            title: "Supplier Name"
        }
    ];
    if (report.isInventory()) {
        columns.push({
            field: "SalesName",
            title: "Sales Name"
        });
    }
    columns.push({
        field: "Total",
        title: "Total",
        template: '#=ChangeToRupiah(FormatCurrency(Total))#',
        attributes: {
            style: "text-align:right;"
        },
    });
    if (report.valueDDReportType2() == "Detail"){
        columns = [
            {
                field:"Date",
                title:"Date",
                template: "#=moment(Date).format('DD-MMM-YYYY')#",
                width:"100px",
            },{
                field:"DocumentNumber",
                title:"Doc. Number"
            },{
                field: "SupplierName",
                title: "Supplier Name"
            }
        ]
        if (report.isInventory())
        {
            columns.push({
                field: "SalesName",
                title: "Sales Name"
            });
        }
        columns.push({
            field: "Item",
            title: "Item"
        });
        columns.push({
            field: "Price",
            title: "Amount",
            template: '#=ChangeToRupiah(FormatCurrency(Price))#',
            attributes: {
                style: "text-align:right;"
            },
        });
        columns.push({
            field: "Qty",
            title: "Qty",
            width:"70px",
            attributes: {
                style: "text-align:center;"
            },
        });
        columns.push({
            field: "Total",
            title: "Total",
            template: '#=ChangeToRupiah(FormatCurrency(Total))#',
            attributes: {
                style: "text-align:right;"
            },
        });
    }
    $("#gridreport").kendoGrid({
        dataSource: {
            data: report.dataGridreport(),
            sort: {
                field: 'Date',
                dir: 'asc',
            }
        },
        excel: {
            fileName: "report-purchaseorder" + (report.isInventory() ? "inventory" : "noninventory") + ".xlsx"
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

report.detailReportPdf = function(){
    var valuefilter = ""
    if (report.valueDDReportType()=="Sales"){
        valuefilter = report.valueSalesCode()
    }else if (report.valueDDReportType()=="Supplier"){
        valuefilter = report.valueSupplierCode()
    }
    var param = {
        DateStart : moment(report.DateStart()).format("YYYY-MM-DD"),
        DateEnd : moment(report.DateEnd()).format("YYYY-MM-DD"),
        ReportType: report.valueDDReportType2(),
        ReportBy : report.valueDDReportType(),
        ValueFilter : valuefilter,
        IsInventory : report.isInventory()
    }
    model.Processing(true)
    ajaxPost("/report/exporttopdfpurchaseorderdetail", param, function(res){
        if (res.IsError){
            model.Processing(false)
            return swal('Error!', res.Mesaage, "errror")
        }
        window.open('/res/docs/report/pdf/' + res, '_blank');
        model.Processing(false)
    })
}
report.ExportToPdf = function(){
    var valuefilter = ""
    if (report.valueDDReportType()=="Sales"){
        valuefilter = report.valueSalesCode()
    }else if (report.valueDDReportType()=="Supplier"){
        valuefilter = report.valueSupplierCode()
    }
    var param = {
        DateStart : moment(report.DateStart()).format("YYYY-MM-DD"),
        DateEnd : moment(report.DateEnd()).format("YYYY-MM-DD"),
        ReportType: report.valueDDReportType2(),
        ReportBy : report.valueDDReportType(),
        ValueFilter : valuefilter,
        IsInventory : report.isInventory()
    }
    model.Processing(true)
    ajaxPost("/report/exporttopdfpurchaseorderreport", param, function(res){
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
report.init = function(){ 
    ProActive.KendoDatePickerRange();    
    report.getDateNow()
    report.getDataSales()
    report.getDataSupplier()
    report.getDataGrid(function(){
        report.renderGrid()
    })
}
report.clearFilter = function() {
    report.valueSalesCode("");
    report.valueSupplierCode("");
    report.getDateNow();
    var data = []
    if (report.isInventory()) {
        data = [
            {
                "Text": "All",
                "Value": "All",
                "TitleReport": "PURCHASE ORDER REPORT"
            },
            {
                "Text": "Sales",
                "Value": "Sales",
                "TitleReport": "PURCHASE ORDER BY SALES"
            },
            {
                "Text": "Supplier",
                "Value": "Supplier",
                "TitleReport": "PURCHASE ORDER BY SUPPLIER"
            }
        ]
    } else {
        data = [
            {
                "Text": "All",
                "Value": "All",
                "TitleReport": "PURCHASE ORDER REPORT"
            },
            {
                "Text": "Supplier",
                "Value": "Supplier",
                "TitleReport": "PURCHASE ORDER BY SUPPLIER"
            }
        ]
    }
    report.dataDDReportType(data);
    report.visibleDDcustomer(false);
    report.visibleDDsales(false);
    report.valueDDReportType("All");
    report.valueDDReportType2("Summary");
}
$(function() {
    report.init()
})