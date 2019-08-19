var activitylog = {}

activitylog.dataMaster = ko.observableArray([])
activitylog.DatePageBar = ko.observable()
activitylog.TitelFilter = ko.observable(" Show Filter");
activitylog.dataYearly = ko.observableArray([])
activitylog.role = ko.observable(false)
activitylog.jumpToDate = new Date()


activitylog.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    activitylog.DatePageBar(page)
}

activitylog.getData = function (callback) {
    var param = {
        DateStart: $('#dateStart').data('kendoDatePicker').value(),
        DateEnd: $('#dateEnd').data('kendoDatePicker').value(),
        Filter: false,
    }
    ajaxPost('/activitylog/getdataactivitylog', param, function (res) {
        if (res.IsError === "true") {
            swal("Error!", res.Message, "error")
            return
        }
        activitylog.dataMaster(res.Data)
        callback()
    }, function () {
        swal("Error!", "Unknown error, please try again", "error")
    })
}

activitylog.renderGrid = function () {
    var data = activitylog.dataMaster();
    if (typeof $('#gridLog').data('kendoGrid') !== 'undefined') {
        $('#gridLog').data('kendoGrid').setDataSource(new kendo.data.DataSource({
            data: data,
            pageSize: 10,
        }))
        return
    }

    var columns = [{
        title: "No.",
        template: function (dataItem) {
            var idxs = _.findIndex(data, function (d) {
                return d._id == dataItem._id
            })

            return idxs + 1
        },
        width: 40,
    }, {
        field: "accesstime",
        title: "Date",
        width: 150,
        template: "#=moment(accesstime).format('DD-MM-YYYY, h:mm:ss a')#"
    }, {
        field: "username",
        title: "User",
        width: 100
    }, {
        field: "activity",
        title: "Activity Log",
        width: 200
    }, {
        field: "desc",
        title: "Description",
        width: 250
    }, ]

    $('#gridLog').kendoGrid({
        dataSource: {
            data: data,
            /*
            transport: {
                read: {
                    url: '/activitylog/getdataactivitylog',
                    dataType: 'json'
                }
            }*/
        },
        height: 300,
        sortable: true,
        scrollable: true,
        columns: columns,
        pageable: {
            refresh: true,
            pageSizes: true,
            buttonCount: 5
        },

    })
}
// function getDates(startDate, stopDate) {
//     var dateArray = [];
//     var currentDate = moment(startDate);
//     var stopDate = moment(stopDate);
//     while (currentDate <= stopDate) {
//         dateArray.push( moment(currentDate).format('DD MMM YYYY') )
//         currentDate = moment(currentDate).add(1, 'days');
//     }
//     return dateArray;
// }
// activitylog.renderGridTimeLineLog = function(){
//     var dateEnd = $('#dateEnd').data('kendoDatePicker').value()
//     var dateStart = new Date(dateEnd)
//     dateStart.setDate(dateStart.getDate() - 7);
//     var range = getDates(dateStart,dateEnd)
//     var dataDate =[]
//     for(i in range){
//         dataDate.push({
//             field : "",
//             title : range[i]
//         })
//     }
//     $('#gridLogTimeLine').kendoGrid({
//         dataSource: {
//             data: activitylog.dataMaster(),
//         },
//         // height: 500,
//         sortable: true,
//         scrollable: true,
//         columns: dataDate,
//         pageable: {
//             refresh: true,
//             pageSizes: true,
//             buttonCount: 5
//         },

//     })
// }
activitylog.deleteActivityLog = function () {
    var data = []
    for (var i = 0; i < activitylog.dataMaster().length; i++) {
        data.push(activitylog.dataMaster()[i]._id)
    }


    model.Processing(true);
    swal({
        title: "Are you sure to delete these data",
        text: "Your will not be able to recover this data",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#DD6B55",
        confirmButtonText: "Yes, do it!",
        cancelButtonText: "No!",
        closeOnConfirm: true,
        closeOnCancel: false
    }, function (isConfirm) {
        if (isConfirm) {
            model.Processing(true)
            var param = {
                Id: data
            }
            var url = '/activitylog/delete'
            ajaxPost(url, param, function (e) {
                if (e.Message == "OK") {

                    setTimeout(function () {
                        swal({
                            title: "Success!",
                            text: "Data has been deleted!",
                            type: "success"
                        }, function () {
                            // location.reload()
                        });
                        // swal("Success!", "Data has been deleted", "success")
                    }, 100)
                } else {
                    swal('Warning', e, 'error');
                    model.Processing(false);
                }
            }, undefined);
        } else {
            swal("Cancelled", "", "error");
            model.Processing(false);
        }
    });
}
activitylog.setDate = function () {
    var datepicker = $("#dateStart").data("kendoDatePicker");
    var now = new Date();
    var years = moment(now).format("YYYY")
    var Month = moment(now).format("MM")
    var newDate = years + "-" + Month + "-" + "01"
    datepicker.value(new Date(newDate))
}
activitylog.renderTimeLine = function(){
    var data = activitylog.dataMaster();
    var userData = _.uniqBy(activitylog.dataMaster(), function (e) {
      return e.username;
    });
    var newUser = []
    var color = ["#f8a398","#51a0ed","#56ca85","#c72222","#ce9904","#10ab0e","#0a748c"]
    for(i in userData){
        newUser.push({
            text: userData[i].username,
            value: userData[i].username,
            color: color[i]
        })
    }
    var newData = []
    for(i in data){
        var dateEnd = new Date(data[i].accesstime)
        dateEnd.setMinutes (dateEnd.getMinutes() + 120);
        newData.push({
            start: new Date(data[i].accesstime),
            end : dateEnd,
            title : data[i].activity,
            description : data[i].desc,
            attendees : [data[i].username]
        })
    }
    console.log(newData)
    $("#LogTimeLine").kendoScheduler({
        date: $("#dateEnd").data("kendoDatePicker").value(),
        workDayStart: new Date("2013/1/1 09:00 AM"),
        workDayEnd: new Date("2013/1/1 6:00 PM"),
        dataSource:{
            data: newData
        },
        views: [ "timeline", "timelineWeek"],
        editable: {
            move: false,
            destroy: false,
            create: false,
            resize: false,
        },
        group: {
            resources: ["Attendees"],
            orientation: "vertical"
        },
        resources: [
            {
                field: "attendees",
                name: "Attendees",
                dataSource: newUser,
                multiple: true,
                title: "Attendees"
            }
        ]
    });
}
activitylog.getDataforCalenderLog = function(callbacks){
    var dateEnd = $("#dateEnd").data("kendoDatePicker").value();
    var years = moment(dateEnd).format("YYYY")
    var dateStart = new Date(years + "-" + "01" + "-" + "01")
    var url = "/activitylog/getdataactivitylogyearly"
    var param = {
        DateStart : dateStart.toISOString(),
        DateEnd : dateEnd.toISOString()
    }
    model.Processing(true)
    ajaxPost(url, param, function(res){
        activitylog.dataYearly(res.Data)
        model.Processing(false)
        callbacks();
    })
}
activitylog.renderCalenderLog = function(){
    $('#calendarLog').fullCalendar({
        header: {
            left: 'today',
            center: 'prev, title ,next',
            right: 'month,agendaWeek,agendaDay,listMonth'
        },
        editable: false,
        events : activitylog.dataYearly(),
        height: 650,
        // put your options and callbacks here
    })
}
activitylog.init = function () {
    activitylog.getData(function () {
        activitylog.renderGrid()
        activitylog.renderTimeLine()
        activitylog.getDateNow()
        // activitylog.renderGridTimeLineLog()
        activitylog.role(userinfo.rolenameh() == "administrator")
    })
    activitylog.getDataforCalenderLog(function(){
        activitylog.renderCalenderLog();
        activitylog.appendJumpTo();
    })
}

activitylog.refresh = function () {
    $("#LogTimeLine").html("")
    activitylog.init()
}

activitylog.appendJumpTo = function() {
    var html = '<span class="fc-button fc-button-jumpto fc-state-default fc-corner-left fc-corner-right">Jump to:</span>';
    html += '<input type="text" id="dateJump" style="width:130px; font-size:12px;">';
    $("span.fc-button-today").after(html);
    $("#dateJump").kendoDatePicker({value: activitylog.jumpToDate});
    $(".fc-button-jumpto").hover(function() {
        $(this).addClass("fc-state-hover");
    }, function() {
        $(this).removeClass("fc-state-hover");
    });
    $(".fc-button-jumpto").click(function() {
        $(this).addClass("fc-state-down");
        setTimeout(function() {
            $(this).removeClass("fc-state-down");
        }.bind(this), 100);

        activitylog.jumpToDate = $("#dateJump").data("kendoDatePicker").value();

        var calendar = $('#calendarLog').data("fullCalendar");
        calendar.gotoDate( activitylog.jumpToDate );
    });
}


$(function () {
    activitylog.setDate()
    activitylog.init()
})