var stock= {}
stock.DatePageBar = ko.observable()
stock.DateStart = ko.observable(new Date)
stock.DateEnd = ko.observable(new Date)
stock.dataGridStock = ko.observableArray([])
stock.dataDropDownItem = ko.observableArray([])
stock.valueDDItem =ko.observable("")
stock.valueLocationSearch = ko.observable(0)
stock.dataDropDownLocation = ko.observableArray([])
stock.getDateNow = function() {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    stock.DatePageBar(page)
}
stock.getdataItemDropdown = function(){
    ajaxPost("/report/getdataitem",{}, function(res){
        if (res.IsError){
            return swal('Error!', res.Mesaage, "errror")
        }
        var data =[]
        _.each(res.Data, function(e){
            data.push({
                text : e.INVID + "-" +e.INVDesc,
                value : e.INVID
            })
        })
        stock.dataDropDownItem(data)
    })
}

stock.getDefaultLocationID = function() {
    var locData = stock.dataDropDownLocation()
    if (locData.length > 0) return locData[0].LocationID;
    return -1;
}

stock.getDatastockingItem = function(callback){
    var param = {
        DateStart : moment(stock.DateStart()).format("YYYY-MM-DD"),
        DateEnd : moment(stock.DateEnd()).format("YYYY-MM-DD"),
        ItemCode : stock.valueDDItem(),
        LocationID : parseInt(stock.valueLocationSearch()) || stock.getDefaultLocationID()
    }
    model.Processing(true)
    ajaxPost("/report/getdatastockitem", param, function(res){
        if (res.IsError){
            model.Processing(false)
            return swal('Error!', res.Mesaage, "errror")
        }
        stock.dataGridStock(res.Data)
        callback();
        model.Processing(false)
    })
}
stock.renderGrid = function(){
    $("#gridstock").html("")
    $("#gridstock").kendoGrid({
        dataSource: {
            data: stock.dataGridStock(),
            sort: {
                field: 'Date',
                dir: 'asc',
            }
        },
        excel: {
            fileName: "report-stockinventory.xlsx"
        },
        height: 500,
        width: 140,
        filterable: true,
        scrollable: true,
        columns: [
            {
                field:"Item",
                title:"Item"
            },{
                field:"StoreLocationName",
                title:"Location"
            },{
                field:"FirstStock",
                title:"Stock"
            },{
                field: "PO",
                title: "PO"
            },{
                field: "CMV",
                title: "CMV"
            },{
                field: "INV",
                title: "INV"
            },{
                field: "CMI",
                title: "CMI"
            },{
                field: "TS",
                title: "TS"
            },{
                field: "TR",
                title: "TR"
            },{
                field: "TotalStock",
                title: "Total Stock"
            }
        ],
    })
}
stock.exportExcel = function(){
    $("#gridstock").getKendoGrid().saveAsExcel();
}
stock.ExportToPdf = function(){
    var param = {
        DateStart : moment(stock.DateStart()).format("YYYY-MM-DD"),
        DateEnd : moment(stock.DateEnd()).format("YYYY-MM-DD"),
        ItemCode : stock.valueDDItem(),
        LocationID : parseInt(stock.valueLocationSearch()) || stock.getDefaultLocationID()
    }
    model.Processing(true)
    ajaxPost("/report/exportstockitempdf", param, function(res){
        if (res.IsError){
            model.Processing(false)
            return swal('Error!', res.Mesaage, "errror")
        }
        window.open('/res/docs/report/pdf/' + res.Data, '_blank');
        model.Processing(false)
    })
}
stock.refreshGrid = function(){
    stock.getDatastockingItem(function(){
        stock.renderGrid()
    })
}
stock.getDataLocation = function(callback) {
    $.ajax({
        url: "/transferorder/getuserlocations",
        success: function (json) {           
            stock.dataDropDownLocation.removeAll();
            stock.dataDropDownLocation(json);
            if (json.length > 0)
            {
                stock.valueLocationSearch(json[0].LocationID);
                stock.valueLocationSearch.valueHasMutated();
            }
            if (typeof callback == "function") callback();
        }
    });
}
stock.onChangeDateStart = function(val){
    if (val.getTime()>stock.DateEnd().getTime()){
        stock.DateEnd(val)
    }
}
stock.init = function(){
    stock.getDateNow()
    stock.getdataItemDropdown()
    stock.getDataLocation(function() {
        stock.getDatastockingItem(function(){
            stock.renderGrid()
        })
    })
}
$(function() {
    stock.init()
})