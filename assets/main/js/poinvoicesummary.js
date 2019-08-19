var summary ={}
summary.DateStart = ko.observable(moment().startOf('month').format('YYYY-MM-DD hh:mm'))
summary.DateEnd = ko.observable(new Date)
summary.koDatePageBar = ko.observable("")
summary.dataDropDownSupplierFilter = ko.observableArray([])
summary.textSupplierSearch = ko.observable()
summary.dataDropDownCustomerFilter = ko.observableArray([])
summary.textCustomerSearch = ko.observable()
summary.filterIndicator = ko.observable(false)
summary.textSearch= ko.observable('')
summary.dataListPO= ko.observableArray([])
summary.dataListPOInventory = ko.observableArray([])
summary.dataListInv= ko.observableArray([])
summary.tabActive = ko.observable(0)

summary.getDateNow = function () {
    var start = moment().startOf('year').format("DD MMMM YYYY")
    var now = moment().startOf('day').format("DD MMMM YYYY")
    var range = start + " to " + now

    var dateFrom = moment(dateFrom).startOf('year').subtract(1, 'year').format("DD MMMM YYYY");
    var dateEnd = moment(dateFrom).endOf('year').subtract(dateFrom, 'year').format("DD MMMM YYYY");
    var rangePrev = dateFrom + " to " + dateEnd + " (Prev)"
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    summary.koDatePageBar(page)
}
summary.getDataSupplier = function() {
    model.Processing(true)
	ajaxPost('/transaction/getsupplier', {}, function(res) {
		if (res.Total === 0) {
		 	swal("Error!", res.Message, "error")
		 	return
		}
		var DataSupplier = res.Data
		for (i in DataSupplier) {
			DataSupplier[i].Kode = DataSupplier[i].Kode + ""
			DataSupplier[i].Name = DataSupplier[i].Kode + " - " + DataSupplier[i].Name
		}
		summary.dataDropDownSupplierFilter(DataSupplier)
		model.Processing(false)
	})
}
summary.setDate = function() {
    var datepicker = $("#dateStartPO").data("kendoDatePicker");
    var now = new Date();
    var years = moment(now).format("YYYY")
    var Month = moment(now).format("MM")
    var newDate = years + "-" + Month + "-" + "01"
	datepicker.value(new Date(newDate))

    var datepicker2 = $("#dateStartInv").data("kendoDatePicker");
    var now2 = new Date();
    var years2 = moment(now).format("YYYY")
    var Month2 = moment(now).format("MM")
    var newDate2 = years + "-" + Month + "-" + "01"
	datepicker2.value(new Date(newDate))
	
	var datepicker3 = $("#dateStartPOINV").data("kendoDatePicker");
    var now = new Date();
    var years = moment(now).format("YYYY")
    var Month = moment(now).format("MM")
    var newDate = years + "-" + Month + "-" + "01"
    datepicker3.value(new Date(newDate))
}
summary.getDataListTrackPO = function(callback){
	var startdate = $('#dateStartPO').data('kendoDatePicker').value();
	var enddate = $('#dateEndPO').data('kendoDatePicker').value();

	var param = {}
	if (summary.filterIndicator() == true) {
		param = {
		    DateStart: startdate,
		    DateEnd: enddate,
		    Filter: true,
		    TextSearch: summary.textSearch(),
		    SupplierCode: summary.textSupplierSearch(),
		}
	} else {
		param = {
			DateStart: startdate,
		    DateEnd: enddate,
		    Filter: false,
		}
	}
	model.Processing(true)
	ajaxPost("/poinvoicesummary/getdatalisttrackpo", param, function(res){
		if ((res.Data).length == 0) {
			res.Data = []
		}
		summary.dataListPO(res.Data)
		model.Processing(false)
		callback();
	});
}
summary.renderGridPO = function(){
	var data = summary.dataListPO()
	// console.log(data)
	$('#gridListPO').kendoGrid({
		dataSource : {
			data: data,
			sort:({
				field:"DateCreated",
				dir: "desc"
			})
		},
		height: 500,
        // width: 140,
        sortable: true,
        scrollable: true,
        columns: [{
         	title: 'Action',
	        width: 70,
	        template: "# if (userinfo.rolenameh() == 'administrator' || userinfo.rolenameh() == 'supervisor' ) {#<button onclick='summary.viewDataPO(\"#: DocumentNumber #\",\"#: Status #\",\"#: SupplierCode #\",\"PO\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button> #}else{#<button onclick='summary.viewDataPO(\"#: DocumentNumber #\",\"#: Status #\",\"#: SupplierCode #\",\"PO\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button> #}#",
	    },{
	        field: 'DateStr',
	        title: 'DateCreated',
	        width: 100,
	    }, {
	        field: 'DocumentNumber',
    	    title: 'Document Number',
	        width: 160,
	    }, {
	        field: 'SupplierName',
	        title: 'Supplier Name',
	        width: 200,
	    },{
	        field: 'DatePO',
	        title: 'Purchase Order',
	        width: 150,
	        template: function(e){
	        	var date = moment(e.DatePO).format("DD-MMM-YYYY")
	        	if(date=="01-Jan-0001"){
	        		date = ""
	        	}
	        	return date
	        }
	    },{
	        field: 'DatePI',
	        title: 'Purchase Invoice',
	        width: 150,
	        template: function(e){
	        	var date = moment(e.DatePI).format("DD-MMM-YYYY")
	        	if(date=="01-Jan-0001"){
	        		date = ""
	        	}
	        	return date
	        }
	    },{
	        field: 'DatePP',
	        title: 'Purchase Payment',
	        width: 150,
	        template: function(e){
	        	var date = moment(e.DatePP).format("DD-MMM-YYYY")
	        	if(date=="01-Jan-0001"){
	        		date = ""
	        	}
	        	return date
	        }
	    },{
	        field: 'Status',
	        title: 'Status',
	        width: 70,
	    },{
	        field: 'Remark',
	        title: 'Remark',
	        width: 200,
	    }],
	    dataBound: function(e){
	    	// console.log(this)
	    	var columns = e.sender.columns;
	    	var PO = this.wrapper.find(".k-grid-header [data-field=" + "DatePO" + "]").index();
	    	var PI = this.wrapper.find(".k-grid-header [data-field=" + "DatePI" + "]").index();
	    	var PP = this.wrapper.find(".k-grid-header [data-field=" + "DatePP" + "]").index();
	    	dataView = this.dataSource.view();
	    	for (var i = 0; i < dataView.length; i++) {
	    		if (dataView[i].Status == "PO") {
	    			var row1 = e.sender.tbody.find("[data-uid='" + dataView[i].uid + "']");
	    			var cell1 = row1.children().eq(PO);
	    			var cell2 = row1.children().eq(PI);
	    			var cell3 = row1.children().eq(PP);
	    			cell1.addClass('greenBackground')
	    			cell2.addClass('greyBackground')
	    			cell3.addClass('greyBackground')
	    		} else if (dataView[i].Status == "PI") {
	    			var row1 = e.sender.tbody.find("[data-uid='" + dataView[i].uid + "']");
	    			var cell1 = row1.children().eq(PI);
	    			var cell2 = row1.children().eq(PP);
	    			cell1.addClass('greenBackground')
	    			cell2.addClass('greyBackground')
	    		} else if (dataView[i].Status == "PP Pending") {
	    			var row1 = e.sender.tbody.find("[data-uid='" + dataView[i].uid + "']");
	    			var cell1 = row1.children().eq(PP);
	    			cell1.addClass('greenBackground')
	    		} else if (dataView[i].Status == "PP Paid") {
	    			var row1 = e.sender.tbody.find("[data-uid='" + dataView[i].uid + "']");
	    			var cell1 = row1.children().eq(PP);
	    			cell1.addClass('greenBackground')
	    		}
	    	}
	    },
        excelExport: function(e) {
            ProActive.kendoExcelRender(e, "PurchaseOrderSummary", function(row, sheet){
                for(var ci = 0; ci < row.cells.length; ci++)
                {
                    var cell = row.cells[ci];
                    if (row.type == "data")
                    {
                        if (ci == 0 || ci == 3 || ci == 4 || ci == 5) {
							cell.format = "dd-mmm-yyyy";
							if (ci == 3 || ci == 4 || ci == 5) {
								cell.value = moment(cell.value).format("DD-MMM-YYYY")
								if (cell.value == "01-Jan-0001") {
									cell.value = ""
								}
							}
                        }
                    }
                }
            });
        },
	})
}

summary.exportExcel = function (typ) {
	if (typ == "PO")
    	$("#gridListPO").getKendoGrid().saveAsExcel();
	if (typ == "POINV")
		$("#gridListPOINV").getKendoGrid().saveAsExcel();
	if (typ == "INV")
		$("#gridListInvoice").getKendoGrid().saveAsExcel();
}

summary.viewDataPO = function(DocNum, Status,SupplierCode, type){
	//console.log(type)
	if (Status == "PO"){
		return window.location.assign("/transaction/purchaseorder?id="+DocNum+"&type="+type)
	}
	if (Status == "PI"){
		return window.location.assign("/transaction/purchaseinvoice?id="+DocNum+"&type="+type)
	}
	if (Status == "PP Pending"){
		return window.location.assign("/transaction/purchasepayment?id="+SupplierCode+"&type="+type)
	}
	if (Status == "PP Paid"){
		// return swal("Warning!", "Document is paid", "warning")
		return window.location.assign("/transaction/purchasepayment?id="+SupplierCode+"&stat=PAID"+"&type="+type)
	}
}
summary.getDataListTrackPOInventory = function(callback){
	var startdate = $('#dateStartPOINV').data('kendoDatePicker').value();
	var enddate = $('#dateEndPOINV').data('kendoDatePicker').value();

	var param = {}
	if (summary.filterIndicator() == true) {
		param = {
		    DateStart: startdate,
		    DateEnd: enddate,
		    Filter: true,
		    TextSearch: summary.textSearch(),
		    SupplierCode: summary.textSupplierSearch(),
		}
	} else {
		param = {
			DateStart: startdate,
		    DateEnd: enddate,
		    Filter: false,
		}
	}
	console.log(param)
	model.Processing(true)
	ajaxPost("/poinvoicesummary/getdatalisttrackpoinventory", param, function(res){
		//console.log(res.Data)
		if ((res.Data).length == 0) {
			res.Data = []
		}
		summary.dataListPOInventory(res.Data)
		model.Processing(false)
		callback();
	});
}
summary.renderGridPOInventory = function(){
	var data = summary.dataListPOInventory()
	// console.log("data")
	// console.log(data)
	$('#gridListPOINV').kendoGrid({
		dataSource : {
			data: data,
			sort:({
				field:"DateCreated",
				dir: "desc"
			})
		},
		height: 500,
        // width: 140,
        sortable: true,
        scrollable: true,
        columns: [{
         	title: 'Action',
	        width: 70,
	        template: "# if (userinfo.rolenameh() == 'administrator' || userinfo.rolenameh() == 'supervisor' ) {#<button onclick='summary.viewDataPO(\"#: DocumentNumber #\",\"#: Status #\",\"#: SupplierCode #\", \"POINV\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button> #}else{#<button onclick='summary.viewDataPO(\"#: DocumentNumber #\",\"#: Status #\",\"#: SupplierCode #\",\"POINV\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button> #}#",
	    },{
	        field: 'DateStr',
	        title: 'DateCreated',
	        width: 100,
	    }, {
	        field: 'DocumentNumber',
    	    title: 'Document Number',
	        width: 160,
	    }, {
	        field: 'SupplierName',
	        title: 'Supplier Name',
	        width: 200,
	    },{
	        field: 'DatePO',
	        title: 'Purchase Order',
	        width: 150,
	        template: function(e){
	        	var date = moment(e.DatePO).format("DD-MMM-YYYY")
	        	if(date=="01-Jan-0001"){
	        		date = ""
	        	}
	        	return date
	        }
	    },{
	        field: 'DatePI',
	        title: 'Purchase Invoice',
	        width: 150,
	        template: function(e){
	        	var date = moment(e.DatePI).format("DD-MMM-YYYY")
	        	if(date=="01-Jan-0001"){
	        		date = ""
	        	}
	        	return date
	        }
	    },{
	        field: 'DatePP',
	        title: 'Purchase Payment',
	        width: 150,
	        template: function(e){
	        	var date = moment(e.DatePP).format("DD-MMM-YYYY")
	        	if(date=="01-Jan-0001"){
	        		date = ""
	        	}
	        	return date
	        }
	    },{
	        field: 'Status',
	        title: 'Status',
	        width: 70,
	    },{
	        field: 'Remark',
	        title: 'Remark',
	        width: 200,
	    }],
	    dataBound: function(e){
	    	// console.log(this)
	    	var columns = e.sender.columns;
	    	var PO = this.wrapper.find(".k-grid-header [data-field=" + "DatePO" + "]").index();
	    	var PI = this.wrapper.find(".k-grid-header [data-field=" + "DatePI" + "]").index();
	    	var PP = this.wrapper.find(".k-grid-header [data-field=" + "DatePP" + "]").index();
	    	dataView = this.dataSource.view();
	    	for (var i = 0; i < dataView.length; i++) {
	    		if (dataView[i].Status == "PO") {
	    			var row1 = e.sender.tbody.find("[data-uid='" + dataView[i].uid + "']");
	    			var cell1 = row1.children().eq(PO);
	    			var cell2 = row1.children().eq(PI);
	    			var cell3 = row1.children().eq(PP);
	    			cell1.addClass('greenBackground')
	    			cell2.addClass('greyBackground')
	    			cell3.addClass('greyBackground')
	    		} else if (dataView[i].Status == "PI") {
	    			var row1 = e.sender.tbody.find("[data-uid='" + dataView[i].uid + "']");
	    			var cell1 = row1.children().eq(PI);
	    			var cell2 = row1.children().eq(PP);
	    			cell1.addClass('greenBackground')
	    			cell2.addClass('greyBackground')
	    		} else if (dataView[i].Status == "PP Pending") {
	    			var row1 = e.sender.tbody.find("[data-uid='" + dataView[i].uid + "']");
	    			var cell1 = row1.children().eq(PP);
	    			cell1.addClass('greenBackground')
	    		} else if (dataView[i].Status == "PP Paid") {
	    			var row1 = e.sender.tbody.find("[data-uid='" + dataView[i].uid + "']");
	    			var cell1 = row1.children().eq(PP);
	    			cell1.addClass('greenBackground')
	    		}
	    	}
	    },
        excelExport: function(e) {
            ProActive.kendoExcelRender(e, "PurchaseOrderInventorySummary", function(row, sheet){
                for(var ci = 0; ci < row.cells.length; ci++)
                {
                    var cell = row.cells[ci];
                    if (row.type == "data")
                    {
                        if (ci == 0 || ci == 3 || ci == 4 || ci == 5) {
							cell.format = "dd-MM-yyyy";
							if (ci == 3 || ci == 4 || ci == 5) {
								cell.value = moment(cell.value).format("DD-MMM-YYYY")
								if (cell.value == "01-Jan-0001") {
									cell.value = ""
								}
							}
                        }
                    }
                }
            });
        },
	})
}
summary.viewDataPOInventory = function(DocNum, Status,SupplierCode){
	if (Status == "PO"){
		return window.location.assign("/transaction/purchaseorder?id="+DocNum)
	}
	if (Status == "PI"){
		return window.location.assign("/transaction/purchaseinvoice?id="+DocNum)
	}
	if (Status == "PP Pending"){
		return window.location.assign("/transaction/purchasepayment?id="+SupplierCode)
	}
	if (Status == "PP Paid"){
		// return swal("Warning!", "Document is paid", "warning")
		return window.location.assign("/transaction/purchasepayment?id="+SupplierCode+"&stat=PAID")
	}
}
summary.getDataCustomer = function () {
    model.Processing(true) 
    ajaxPost('/transaction/getcustomer', {}, function (res) {

        if (res.Total === 0) {
            swal("Error!", res.Message, "error")
            return
        }
        var DataCustomer = res.Data
        for (i in DataCustomer) {
            DataCustomer[i].Kode = DataCustomer[i].Kode + ""
            DataCustomer[i].Name = DataCustomer[i].Kode + "-" + DataCustomer[i].Name
        }
        summary.dataDropDownCustomerFilter(DataCustomer)
        model.Processing(false)
    })
}

summary.getDataListTrackInv = function(callback){
	var startdate = $('#dateStartInv').data('kendoDatePicker').value();
	var enddate = $('#dateEndInv').data('kendoDatePicker').value();

	var param = {}
	if (summary.filterIndicator() == true) {
		param = {
		    DateStart: startdate,
		    DateEnd: enddate,
		    Filter: true,
		    TextSearch: summary.textSearch(),
		    CustomerCode: summary.textCustomerSearch(),
		}
	} else {
		param = {
			DateStart: startdate,
		    DateEnd: enddate,
		    Filter: false,
		}
	}
	model.Processing(true)
	ajaxPost("/poinvoicesummary/getdatalisttrackinv", param, function(res){
		if ((res.Data).length == 0) {
			res.Data = []
		}
		summary.dataListPO(res.Data)
		model.Processing(false)
		callback();
	});
}
summary.renderGridInv = function(){
	var data = summary.dataListPO()
	$('#gridListInvoice').kendoGrid({
		dataSource : {
			data: data,
			sort:({
				field:"DateCreated",
				dir: "desc"
			})
		},
		height: 500,
        // width: 140,
        sortable: true,
        scrollable: true,
        columns: [{
         	title: 'Action',
	        width: 100,
	        template: "# if (userinfo.usernameh() || userinfo.usernameh() == 'administrator' || userinfo.rolenameh() == 'supervisor' ) {#<button onclick='summary.viewDataInv(\"#: DocumentNumber #\",\"#: Status #\",\"#: CustomerCode #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button> #}else{#<button onclick='summary.viewDataInv(\"#: DocumentNumber #\",\"#: Status #\",\"#: CustomerCode #\")' class='btn btn-sm btn-default btn-flat'><i class='fa fa-address-card-o' aria-hidden='true'></i></button> #}#",
	    },{
	        field: 'DateStr',
	        title: 'DateCreated',
	        width: 160,
	        template: function(e){
	        	var date = moment(e.DateINV).format("DD-MMM-YYYY")
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
	        field: 'CustomerName',
	        title: 'Customer Name',
	        width: 200,
	    },{
	        field: 'DateINV',
	        title: 'Invoice',
	        width: 150,
	        template: function(e){
	        	var date = moment(e.DateINV).format("DD-MMM-YYYY")
	        	if(date=="01-Jan-0001"){
	        		date = ""
	        	}
	        	return date
	        }
	    },{
	        field: 'DateSP',
	        title: 'Sales Payment',
	        width: 150,
	        template: function(e){
	        	var date = moment(e.DateSP).format("DD-MMM-YYYY")
	        	if(date=="01-Jan-0001"){
	        		date = ""
	        	}
	        	return date
	        }
	    },{
	        field: 'Status',
	        title: 'Status',
	        width: 200,
	    },{
	        field: 'Remark',
	        title: 'Remark',
	        width: 200,
	    }],
	    dataBound: function(e){
	    	var columns = e.sender.columns;
	    	var INV = this.wrapper.find(".k-grid-header [data-field=" + "DateINV" + "]").index();
	    	var SP = this.wrapper.find(".k-grid-header [data-field=" + "DateSP" + "]").index();
	    	dataView = this.dataSource.view();
	    	for (var i = 0; i < dataView.length; i++) {
	    		if (dataView[i].Status == "INVOICE") {
	    			var row1 = e.sender.tbody.find("[data-uid='" + dataView[i].uid + "']");
	    			var cell1 = row1.children().eq(INV);
	    			var cell2 = row1.children().eq(SP);
	    			cell1.addClass('greenBackground')
	    			cell2.addClass('greyBackground')
	    		} else if (dataView[i].Status == "SP Pending") {
	    			var row1 = e.sender.tbody.find("[data-uid='" + dataView[i].uid + "']");
	    			var cell1 = row1.children().eq(SP);
	    			cell1.addClass('greenBackground')
	    		} else if (dataView[i].Status == "SP Paid") {
	    			var row1 = e.sender.tbody.find("[data-uid='" + dataView[i].uid + "']");
	    			var cell1 = row1.children().eq(SP);
	    			cell1.addClass('greenBackground')
	    		}
	    	}
	    },
        excelExport: function(e) {
            ProActive.kendoExcelRender(e, "PurchaseInvoiceSummary", function(row, sheet){
                for(var ci = 0; ci < row.cells.length; ci++)
                {
                    var cell = row.cells[ci];
                    if (row.type == "data")
                    {
                        if (ci == 0 || ci == 3 || ci == 4) {
							cell.format = "dd-MM-yyyy";
							if (ci == 3 || ci == 4) {
								cell.value = moment(cell.value).format("DD-MMM-YYYY")
								if (cell.value == "01-Jan-0001") {
									cell.value = ""
								}
							}
                        }
                    }
                }
            });
        },
	})
}
summary.search= function(type){
	if (type=="PO"){
		summary.filterIndicator(true)
		summary.getDataListTrackPO(function(){
			summary.renderGridPO()
		});
	}else if(type=="INV"){
		summary.filterIndicator(true)
		summary.getDataListTrackInv(function(){
			summary.renderGridInv()
		});
	}else if (type == "POINV"){
		console.log("hallo ini PO INV")
		summary.filterIndicator(true)
		summary.getDataListTrackPOInventory(function(){
			summary.renderGridPOInventory()
		});
	} else if (type == undefined){
		if (summary.tabActive()==0){
			summary.filterIndicator(true)
			summary.getDataListTrackPO(function(){
				summary.renderGridPO()
			});
		} else if (summary.tabActive()==1){
			summary.filterIndicator(true)
			summary.getDataListTrackPOInventory(function(){
				summary.renderGridPOInventory()
			});
		} else if(summary.tabActive()==2){
			summary.filterIndicator(true)
			summary.getDataListTrackInv(function(){
				summary.renderGridInv()
			});
		}
	} else{
		summary.filterIndicator(true)
		summary.init()
	}
}
summary.viewDataInv = function(DocNum, Status,CustomerCode){
	if (Status == "INVOICE"||Status == "POSTING"){
		return window.location.assign("/transaction/invoice?id="+DocNum)
	}
	if (Status == "SP Pending"){
		return window.location.assign("/transaction/salespayment?id="+CustomerCode)
	}
	if (Status == "SP Paid"){
		// return swal("Warning!", "Document is paid", "warning")
		return window.location.assign("/transaction/salespayment?id="+CustomerCode+"&stat=PAID")
	}
}

summary.changeTab = function(index){
	if (index == 0) {
		summary.tabActive(0)
		summary.textSearch("")
		$("#textSearch").val("")
	} else if(index==1){
		summary.tabActive(1)
		summary.textSearch("")
		$("#textSearch").val("")
	}else {
		summary.tabActive(2)
		summary.textSearch("")
		$("#textSearch").val("")
	}
}

summary.onChangeDateStart = function(val){
    if (val.getTime()>summary.DateEnd().getTime()){
        summary.DateEnd(val)
    }
}

summary.init = function(){
	summary.getDateNow()
	// summary.setDate()
	summary.getDataSupplier()
	summary.getDataCustomer()
	summary.getDataListTrackPO(function(){
		summary.renderGridPO()
	});
	summary.getDataListTrackInv(function(){
		summary.renderGridInv()
	});
	summary.getDataListTrackPOInventory(function(){
		summary.renderGridPOInventory()
	});
}
$(function () {
	summary.init()
})