var locationmaster = {}

locationmaster.DatePageBar = ko.observable()
locationmaster.text = ko.observable("Save")
locationmaster.filterindicator = ko.observable(false)
locationmaster.dataMasterlocationmaster = ko.observableArray([])
locationmaster.dataMasterlocationmasterOriginal = ko.observableArray([])
locationmaster.getDateNow = function(){
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    locationmaster.DatePageBar(page)
}
locationmaster.newRecord = function() {
    var page = {
        ID: "",
        _id: "",
        Main_LocationID:0,
        LocationID: 0,
        LocationName: "",
        Description: "",
    }
    return page
}

locationmaster.record = ko.mapping.fromJS(locationmaster.newRecord())

locationmaster.getDatalocationmaster = function(callback) {
    model.Processing(true)
    ajaxPost('/master/getdatalocation', {}, function (res) {
        model.Processing(false)
        if (res.Total === 0) {
            swal({
                title:"Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        //console.log(res)
        locationmaster.dataMasterlocationmaster(res.Data)
        locationmaster.dataMasterlocationmasterOriginal(res.Data)
        callback()
    }, function () {
        swal({
            title:"Error!",
            text: "Unknown error, please try again",
            type: "error",
            confirmButtonColor:"#3da09a"
        })
    })
}


locationmaster.reset = function() {
    $(".formInput").val("")
    $("#qty").val(0)
}

locationmaster.addNewModal = function() {
    ko.mapping.fromJS(locationmaster.newRecord(), locationmaster.record);
    locationmaster.reset()
    locationmaster.text("Save")
}

locationmaster.renderGrid = function() {
    var data = locationmaster.dataMasterlocationmaster();
    //console.log(data)
    if (typeof $('#gridlocationmaster').data('kendoGrid') !== 'undefined') {
        $('#gridlocationmaster').data('kendoGrid').setDataSource(new kendo.data.DataSource({
            data: data,
        }))
        return
    } 

    var columns = [{
        title: "Action",
        width: 90,
        template: "<button onclick='locationmaster.edit(\"#: _id #\")' class='btn btn-xs btn-warning'><i class='fa fa-pencil'></i></button> <button onclick='locationmaster.delete(\"#: _id #\",\"#: LocationID #\")' class='btn btn-xs btn-danger'><i class='glyphicon glyphicon-trash' aria-hidden='true'></i></button>"
    },{
        field: "LocationID",
        title: "Location ID",
        width: 160,
    },{ 
        field: "LocationName",
        title: "Location Name",
        width: 160,
    },{
        field: "Description",
        title: "Description",
        width: 60,
    }]
    $('#gridlocationmaster').kendoGrid({
        dataSource: {
            data: data,
        },
        height: 500,
        width: 140,
        sortable: true,
        scrollable: true,
        columns: columns,
        excelExport: function(e) {
            ProActive.kendoExcelRender(e, "Location", function(row, sheet){
                
            });
        },
    })
}
locationmaster.exportExcel = function () {
    $("#gridlocationmaster").getKendoGrid().saveAsExcel();
}

locationmaster.saveData = function() {
    var change = ko.mapping.toJS(locationmaster.record)  
    change.ID = change._id || change.ID; 
    change.Main_LocationID = parseInt(change.Main_LocationID) 
    change.LocationID = parseInt(change.LocationID) 
    change.LocationName = change.LocationName.trim().toUpperCase()
    var param = {
        Data: change,
    }

    var name = change.LocationName
    if (name == "") {
        return swal({
            title:'Warning!', 
            text:"You haven't fill the Location Name", 
            type:"info",
            confirmButtonColor:"#3da09a"
        })
    }
    if (change.Main_LocationID < 1000 || change.Main_LocationID >= 10000) {
        return swal({
            title:'Warning!', 
            text:"Main Location ID must be a 4 digits number", 
            type:"info",
            confirmButtonColor:"#3da09a"
        })
    }
    if (change.LocationID < 1000 || change.LocationID >= 10000) {
        return swal({
            title:'Warning!', 
            text:"Location ID must be a 4 digits number", 
            type:"info",
            confirmButtonColor:"#3da09a"
        })
    }

    // checking parent
    var lid = change.LocationID;
    var pid = change.Main_LocationID;
    var mli = null;
    if (lid % 1000 == 0) {
        if (pid != 1000) {
            mli = 1000;
        }
    }
    else if (lid % 100 == 0) {
        if (pid != (lid - lid%1000)) {
            mli = lid - lid%1000;
        }
    }
    else {
        if (pid != (lid - lid%100)) {
            mli = lid - lid%100;
        }
    }
    if (mli != null)
    {
        return swal({
            title:'Warning!', 
            text:"Main Location ID should be " + mli, 
            type:"info",
            confirmButtonColor:"#3da09a"
        })
    }

    var url = "/master/insertnewlocation"
    swal({
        title: "Are you sure?",
        text: "You will submit this Location",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3da09a",
        confirmButtonText: "Yes, do it!",
        cancelButtonText: "No!",
        closeOnConfirm: true,
        closeOnCancel: false
    }, function(isConfirm) {
        if (isConfirm) {
           model.Processing(true)
            ajaxPost(url, param, function(e) {
                if (!e.status) {
                    setTimeout(function () {
                        swal({
                            title: "Warning!",
                            text: e.Message,
                            type: "error",
                            confirmButtonColor:"#3da09a"
                        }, function () {
                            model.Processing(false);
                        });
                    }, 200)
                } else {
                    swal({
                    title:'Success', 
                    text:'Location has been saved!',
                    type:'success',
                    confirmButtonColor:"#3da09a"
                });
                    locationmaster.init()
                    locationmaster.reset()
                    $('#AddNewModal').modal('hide');
                    model.Processing(false);
                }
                model.Processing(false)
            })
        } else {
            swal({
                title:"Cancelled",
                type:"error",
                confirmButtonColor:"#3da09a"
            });
        }
    });
}

locationmaster.edit = function(id) {    
    var data = _.find(locationmaster.dataMasterlocationmaster(), function(o) {
        return o._id == id;
    });
    // need to map to blank record to prevent blank fields when editing data that is previously updated
    ko.mapping.fromJS(locationmaster.newRecord(), locationmaster.record);
    ko.mapping.fromJS(data, locationmaster.record);
    $("#AddNewModal").modal("show")
    locationmaster.text("Update")
}

locationmaster.delete = function(id, INVID) {
    model.Processing(true);
    swal({
        title: "Are you sure to delete " + INVID + "?",
        text: "Your will not be able to recover this data",
        type: "warning",
        showCancelButton: true,
        confirmButtonClass: "#3da09a",
        confirmButtonText: "Yes, delete it!",
        closeOnConfirm: false
    }, function (res) {
        if (res) {
            var url = "/master/deletelocation";
            var param = {
                ID: id
            };
            ajaxPost(url, param, function (data) {
                if (data.Message != "OK") {
                    swal({
                        title:'Warning',
                        text: data.Message,
                        type: 'error',
                        confirmButtonColor:"#3da09a"
                    });
                    model.Processing(false);
                } else {
                    swal({
                        title:'Success',
                        text: 'locationmaster has been deleted!',
                        type: 'success',
                        confirmButtonColor:"#3da09a"
                    });
                    locationmaster.init()
                    model.Processing(false);
                }

            }, undefined);
        } else {
            model.Processing(false);
        }
    });
}

locationmaster.init = function() {
    locationmaster.getDatalocationmaster(function(){
        locationmaster.renderGrid() 
        locationmaster.getDateNow() 
    })
}

locationmaster.filterText = function(term) {
    var filter = { logic: "or", filters: [] };
    var filteredFields = ["LocationName", "Description"]
    $searchValue = term || $("#textSearch").val();
    if ($searchValue) {
        for (var k in filteredFields)
            filter.filters.push({ field: filteredFields[k], operator:"contains", value:$searchValue});
    } 
    $("#gridlocationmaster").data("kendoGrid").dataSource.query({ filter: filter });
}

$(function(){
    locationmaster.init()
    $("#textSearch").on("keyup blur change", function () {
        locationmaster.filterText();
    });
})