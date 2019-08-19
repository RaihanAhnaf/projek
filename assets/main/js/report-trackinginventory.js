var track= {}
track.DatePageBar = ko.observable()
track.DateStart = ko.observable(moment().startOf('month').format('YYYY-MM-DD hh:mm'))
track.DateEnd = ko.observable(new Date)
track.dataGridTrack = ko.observableArray([])
track.dataDropDownItem = ko.observableArray([])
track.datamasterInv = ko.observableArray([])
track.dataMasterLocation = ko.observableArray([])
track.dataDropDownLocation = ko.observableArray([])
track.valueLocation = ko.observable('')
track.valueDDItem =ko.observable("")
track.getDateNow = function() {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    track.DatePageBar(page)
}
track.getdataItemDropdown = function(){
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
        track.datamasterInv(res.Data)
        track.dataDropDownItem(data)
    })
}
track.getdataLocation = function(){
    model.Processing(true)
    ajaxPost("/report/getdatalocationtracking", {}, function(res){
        if (res.IsError){
            model.Processing(false)
            return swal('Error!', res.Mesaage, "errror")
        }
        track.dataMasterLocation(res.Data)
        var data =[]
        _.each(res.Data, function(e){
            data.push({
                text : e.LocationID + "-" +e.LocationName,
                value : e.LocationID
            })
        })
        track.dataDropDownLocation(data)
        track.valueLocation(userinfo.locationid())
        model.Processing(false)
    })
}
track.getDataTrackingItem = function(callback){
    if (track.valueDDItem()==""){
        return swal("Warning", "Please Select One of Item", "info")
    }
    var locationvalue = 0
    if(track.valueLocation() != ""){
        locationvalue = parseInt(track.valueLocation())
    }
    var param = {
        DateStart : moment(track.DateStart()).format("YYYY-MM-DD"),
        DateEnd : moment(track.DateEnd()).format("YYYY-MM-DD"),
        ItemCode : track.valueDDItem(),
        LocationCode : locationvalue,
    }
    model.Processing(true)
    ajaxPost("/report/getdatatrackingitem", param, function(res){
        if (res.IsError){
            model.Processing(false)
            return swal('Error!', res.Mesaage, "errror")
        }
        track.dataGridTrack(res.Data)
        callback();
        model.Processing(false)
    })
}

track.renderGrid = function(){
    $("#gridtrack").html("")
    $("#gridtrack").kendoGrid({
        dataSource: {
            data: track.dataGridTrack(),
            sort: {
                field: 'Id',
                dir: 'asc',
            }
        },
        excel: {
            fileName: "report-trackinventory.xlsx"
        },
        height: 500,
        width: 140,
        filterable: true,
        scrollable: true,
        columns: [
            {
                field : "Date",
                title : "Date",
                template: function(e){
                    return moment(e.Date).format("DD MMM YYYY")
                }
            },{
                field:"Description",
                title:"Description"
            },{
                field:"StoreHouseName",
                title:"Location"
            },{
                field:"TypeTransaction",
                title:"Type"
            },{
                field: "StockUnit",
                title: "Stock"
            },{
                field: "Increment",
                title: "Increment"
            },{
                field: "Decrement",
                title: "Decrement"
            },{
                field: "TotalSaldo",
                title: "Total"
            }
        ],
    })
}
track.exportExcel = function(){
    $("#gridtrack").getKendoGrid().saveAsExcel();
}
track.ExportToPdf = function(){
    if (track.valueDDItem()==""){
        return swal("Warning", "Please Select One of Item", "info")
    }
    var item = _.find(track.datamasterInv(), function(e){return e.INVID== track.valueDDItem()})
    var locationvalue = 0
    if(track.valueLocation() != ""){
        locationvalue = parseInt(track.valueLocation())
    }
    var param = {
        DateStart : moment(track.DateStart()).format("YYYY-MM-DD"),
        DateEnd : moment(track.DateEnd()).format("YYYY-MM-DD"),
        ItemCode : track.valueDDItem(),
        ItemName : item.INVDesc,
        LocationCode : locationvalue,
    }
    model.Processing(true)
    ajaxPost("/report/exporttrackingitemtopdf", param, function(res){
        if (res.IsError){
            model.Processing(false)
            return swal('Error!', res.Mesaage, "errror")
        }
        window.open('/res/docs/report/pdf/' + res.Data, '_blank');
        model.Processing(false)
    })
}
track.onChangeDateStart = function(val){
    if (val.getTime()>track.DateEnd().getTime()){
        track.DateEnd(val)
    }
}
track.refreshGrid = function(){
    track.getDataTrackingItem(function(){
        track.renderGrid()
    })
}
track.init = function(){
    track.getDateNow()
    track.getdataItemDropdown()
    track.getdataLocation()
    track.renderGrid()
}
$(function() {
    track.init()
})