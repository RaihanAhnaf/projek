var report= {}
report.DatePageBar = ko.observable()
report.DateStart = ko.observable(moment().startOf('month').format('YYYY-MM-DD'))
report.DateEnd = ko.observable(new Date)
report.dataTransferShipment = ko.observableArray([]),
report.gridDataSource = ko.observableArray([]),
report.dataGridreport = ko.observableArray([])
report.dataMasterSales = ko.observableArray([])
report.dataDropDownSales= ko.observableArray([])
report.valueSalesCode = ko.observable()
report.visibleDDsales = ko.observable(false)
report.dataMasterCustomer = ko.observableArray([])
report.dataDropDownCustomerFilter = ko.observableArray([])
report.valueCustomerCode = ko.observable()
report.visibleDDcustomer = ko.observable(false)
report.warehouse = ko.observableArray([])
report.valueStorehouse = ko.observable()
report.dataLocation = ko.observableArray([])
report.storeHouseLabel = ko.observable("Store House From :")
report.valueDDReportType = ko.observable('All')
report.dataDDReportType2 = ko.observableArray([
    {
        "Text": "Transfer Shipment",
        "Value": "Transfer Shipment",
    },
    {
        "Text": "Transfer Receipt",
        "Value": "Transfer Receipt",
    }
])
report.valueDDReportType2 = ko.observable('Transfer Shipment')
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

report.onChangeSalesCode = function(value){
    if (value== "Transfer Shipment"){
        report.storeHouseLabel("Store House From :")
    }else if (value== "Transfer Receipt"){
        report.storeHouseLabel("Store House To :")
        // report.visibleDDsales(true)
        // report.visibleDDcustomer(false)
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

report.refreshDataSource = function(filter, callback) {
    filter = filter || {
        Filter: false
    };

    if ($("#storehouse").data("kendoDropDownList").value() != ""){
        filterLocation = $("#storehouse").data("kendoDropDownList").value();
    } else{
        filterLocation = report.valueStorehouse();
    }

    var param = {
        DateStart : moment(report.DateStart()).format("YYYY-MM-DD"),
        DateEnd : moment(report.DateEnd()).format("YYYY-MM-DD"),
        Location : parseInt(filterLocation),
        Filter : true
    }

    var url ="";
    if (report.valueDDReportType2() == "Transfer Shipment"){
         url = '/transferorder/getshipmentlist';
    } else if (report.valueDDReportType2() == "Transfer Receipt"){
         url = '/transferorder/getreceiptlist';
    }
    model.Processing(true)
    ajaxPost(url, param, function(res) {
        if (res.IsError) {
            model.Processing(false)
            swal({
                title: "Search Not Found!",
                text: res.Message,
                type: "warning",
                confirmButtonColor: "#3da09a"
            })
            $('#textSearch').val("")
            return
        }
        model.Processing(false)
        report.dataTransferShipment(res.Data)
        report.gridDataSource.removeAll();
        $(res.Data).each(function(idx, ele) {
            ele.ListDetailOrder = ele.ListDetailTransferShipment;
            delete ele.ListDetailTransferShipment;
            report.gridDataSource.push(ele);
        });

        // if (filter.Filter) {
        //     var grid = $("#gridListTransferShipment").data("kendoGrid");
        //     grid.dataSource.read()
        //     grid.refresh();
        //     //console.log("Grid Refreshed: ", res.Data);
        // }

        model.Processing(false);
        if (callback) callback.apply(report);
    }, function() {
        model.Processing(false)
        swal({
            title: "Error!",
            text: "Unknown error, please try again",
            type: "error",
            confirmButtonColor: "#3da09a"
        })
    })
}

report.renderGrid = function(){
    report.refreshDataSource(null, function() {
        // if (typeof $('#gridListTransferShipment').data('kendoGrid') !== 'undefined') {
        //     console.log("tes111")
        //     $('#gridListTransferShipment').data('kendoGrid').setDataSource(new kendo.data.DataSource({
        //         data: report.gridDataSource(),
        //     }))
        //     return
        // }

        if (report.valueDDReportType2() == "Transfer Shipment"){
            DocumentNumTitle = 'Document Number Shipment';
            $("#gridListTransferShipment th[data-field=DocumentNumberShipment]").contents().last().replaceWith("Document Number Shipment")
            DocumentNumField = 'DocumentNumberShipment';
            DetailButton = "<button onclick='report.viewDataTO(\"#: DocumentNumberShipment #\",\"TS\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button>";
       } else if (report.valueDDReportType2() == "Transfer Receipt"){
            DocumentNumTitle = 'Document Number Receipt';
            $("#gridListTransferShipment th[data-field=DocumentNumberReceipt]").contents().last().replaceWith("Document Number Receipt")
            DocumentNumField = 'DocumentNumberReceipt';
            DetailButton = "<button onclick='report.viewDataTO(\"#: DocumentNumberReceipt #\",\"TR\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button>";
       }

        var columns = [{
            title: 'Action',
            width: 100,
            template: DetailButton,
        }, {
            field: DocumentNumField,
            title: DocumentNumTitle,
            width: 200,
        }, {
            field: 'StoreHouseFrom', title: 'Store House From', width: 200,
            template : function(dt){
                return dt.FromDetail[0].LocationName + " (" + dt.StoreHouseFrom + ")"
            }
        }, {
            field: 'StoreHouseTo', title: 'Store House To', width: 200,
            template : function(dt){
                return dt.ToDetail[0].LocationName + " (" + dt.StoreHouseTo + ")"
            }
        }, {
            field: 'DateStr', title: 'Date', width: 110,
        }, {
            field: 'Description', title: 'Description', width: 200,
        }]

        $('#gridListTransferShipment').kendoGrid({
            dataSource: {
                data: report.gridDataSource(),
                sort: {
                    field: 'DateStr',
                    dir: 'desc',
                }
            },
            height: 500,
            width: 140,
            sortable: true,
            scrollable: true,
            columns: columns,
            excelExport: function(e) {
                titleFile = "report " +  report.valueDDReportType2() + " " + moment(report.DateStart()).format("YYYYMMDD") + " to "+ moment(report.DateEnd()).format("YYYYMMDD") + " " + $("#storehouse").data("kendoDropDownList").text()
                ProActive.kendoExcelRender(e, titleFile, function(row, sheet){
                    for(var ci = 0; ci < row.cells.length; ci++)
                    {
                        var cell = row.cells[ci];
                        if (row.type == "data")
                        {
                            if (ci == 1 || ci == 2)
                                cell.hAlign = "left";
                            if (ci == 3)
                                cell.format = "dd-MM-yyyy";
                        }
                    }
                });
            },
        });

        if (report.valueDDReportType2() == "Transfer Shipment"){
            $("#gridListTransferShipment th[data-field=DocumentNumberShipment]").contents().last().replaceWith("Document Number Shipment")
       } else if (report.valueDDReportType2() == "Transfer Receipt"){
            $("#gridListTransferShipment th[data-field=DocumentNumberReceipt]").contents().last().replaceWith("Document Number Receipt")
       }
    });
}

report.exportExcel = function(){
    $("#gridListTransferShipment").getKendoGrid().saveAsExcel();
}

report.ExportToPdf = function(){
    var locationData = _.find(report.warehouse(), function(o) {
        return o.value == report.valueStorehouse();
    });
    
    var param = {
        DateStart : moment(report.DateStart()).format("YYYY-MM-DD"),
        DateEnd : moment(report.DateEnd()).format("YYYY-MM-DD"),
        TransferType: report.valueDDReportType2(),
        StoreHouse : locationData.text +" ("+ report.valueStorehouse()+")",
        StoreHouseVal : parseInt(report.valueStorehouse())
    }
    model.Processing(true)
    ajaxPost("/transferorder/exporttopdfreporttransferorder", param, function(res){
        if (res.IsError){
            model.Processing(false)
            return swal('Error!', res.Mesaage, "errror")
        }
        window.open('/res/docs/report/pdf/' + res.Data, '_blank');
        model.Processing(false)
    })
}

report.detailReportPdf = function(){
    var locationData = _.find(report.warehouse(), function(o) {
        return o.value == report.valueStorehouse();
    });

    if (report.valueDDReportType2() == "Transfer Shipment") {
        url = "/transferorder/exportdetailreportallts"
    } else if (report.valueDDReportType2() == "Transfer Receipt"){
        url = "/transferorder/exportdetailreportalltr"
    }
    
    var param = {
        DateStart : moment(report.DateStart()).format("YYYY-MM-DD"),
        DateEnd : moment(report.DateEnd()).format("YYYY-MM-DD"),
        TransferType: report.valueDDReportType2(),
        StoreHouse : locationData.text +" ("+ report.valueStorehouse()+")",
        StoreHouseVal : parseInt(report.valueStorehouse())
    }
    model.Processing(true)
    ajaxPost(url, param, function(res){
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

report.getMasterLocation= function() {
    $.ajax({
        url: "/transferorder/getuserlocations",
        success: function(data) {
            report.warehouse.removeAll();
            $(data).each(function(ix, ele) {
                report.dataLocation.push({
                    value: ele.LocationID,
                    text: ele.LocationName
                });
                report.warehouse.push({
                    value: ele.LocationID,
                    text: ele.LocationName
                });
            });
            report.valueStorehouse(data[0]["LocationID"])
        }
    });
}

report.filterText = function(term) {
    var filter = { logic: "or", filters: [] };
    var filteredFields = ["DocumentNumberShipment", "DocumentNumberReceipt", "Description"]
    $searchValue = term || $("#textSearch").val();
    if ($searchValue) {
        for (var k in filteredFields)
            filter.filters.push({ field: filteredFields[k], operator:"contains", value:$searchValue});
    } 
    $("#gridListTransferShipment").data("kendoGrid").dataSource.query({ filter: filter });
}

report.init = function(){
    report.getDateNow()
    report.getMasterLocation()
    report.initDate();
    report.getDataGrid(function(){
        report.renderGrid()
    })
}
report.initDate = function () {
    var dtpStart = $('#dateStart').data('kendoDatePicker');
    var dtpEnd = $('#dateEnd').data('kendoDatePicker');
    dtpStart.value(moment().startOf('month').toDate());
    dtpEnd.value(moment().startOf('day').toDate());
    dtpStart.max(dtpEnd.value());
    dtpEnd.min(dtpStart.value());

    dtpStart.bind("change", function () {
        dtpEnd.min(dtpStart.value());
    });
    dtpEnd.bind("change", function () {
        dtpStart.max(dtpEnd.value());
    });
}

report.viewDataTO = function(DocNum, type){
	if (type == "TS"){
		return window.location.assign("/transferorder/transfershipment?id="+DocNum+"&type="+type)
	}
	if (type == "TR"){
		return window.location.assign("/transferorder/transferreceipt?id="+DocNum+"&type="+type)
	}
}
$(function() {
    report.init()
    $("#textSearch").on("keyup blur change", function () {
        report.filterText();
    });
})