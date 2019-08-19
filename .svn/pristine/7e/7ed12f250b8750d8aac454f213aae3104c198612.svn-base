var asset = {}

asset.dataMasterAsset = ko.observableArray([])
asset.dataMasterAssetOriginal = ko.observableArray([])
asset.dataMasterCategory = ko.observableArray([])
asset.DatePageBar = ko.observable()
asset.button = ko.observable("Save")
asset.dataFilter = ko.observableArray([])
asset.textTitle = ko.observable("Add New Asset")
asset.textSearch = ko.observable()
asset.filterindicator = ko.observable(false)
asset.dataDepreciationMonth = ko.observableArray([])
asset.dataIdcheckBox = ko.observableArray([])
asset.dataModalDepreciation = ko.observableArray([])
asset.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    asset.DatePageBar(page)
}

asset.newRecord = function () {
    return {
        ID: ko.observable(),
        Description: ko.observable(''),
        Category: ko.observable(''),
        Qty: ko.observable(0),
        Price: ko.observable(0),
        Total: ko.observable(0),
        PostingDate: ko.observable(''),
        DatePeriod: ko.observable(''),
        SumDepreciation: ko.observable(0),
        MonthlyDepreciation: ko.observable(0),
        User: ko.observable(''),
    }
}
asset.record = ko.mapping.fromJS(asset.newRecord())

asset.record.Total = ko.computed(function () {
    var n = String(asset.record.Price())
    var total = asset.record.Qty() * parseFloat(n.replace(/,/g, ''))
    return total.toFixed(2);
})

asset.record.SumDepreciation = ko.computed(function () {
    var month = 0
    var a = moment(asset.record.PostingDate())
    var b = moment(asset.record.DatePeriod())
    month = b.diff(a, 'month')
    return (month + 1);
})

asset.record.MonthlyDepreciation = ko.computed(function () {
    var month = 0
    var a = moment(asset.record.PostingDate())
    var b = moment(asset.record.DatePeriod())
    month = b.diff(a, 'month')
    var total = asset.record.Total() / (month + 1)
    return total.toFixed(2);
})

asset.saveData = function () {
    var validator = $("#myForm").data("kendoValidator");
    if (validator === undefined) {
        validator = $("#myForm").kendoValidator().data("kendoValidator");
    }
    if (!validator.validate()) {
        return
    }
    var assetPrice = FormatCurrency(asset.record.Price())

    var url = "/transaction/insertnewasset"

    var payload = ko.mapping.toJS(asset.record)
    payload.MonthlyDepreciation = FormatCurrency(asset.record.MonthlyDepreciation())
    payload.Price = assetPrice
    payload.Qty = FormatCurrency(asset.record.Qty())
    payload.Total = FormatCurrency(asset.record.Total())
    payload.User = userinfo.usernameh()
    payload.PostingDate = $('#postingdate').data('kendoDatePicker').value()
    payload.DatePeriod = $('#dateperiod').data('kendoDatePicker').value()

    var param = {
        Data: JSON.stringify(payload) 
        // Data : payload      
    }
    model.Processing(true)
    ajaxPost(url, param, function (res) {
        if (res.IsError === "true") {
            swal({
            title:"Error!",
            text: res.Message,
            type: "error",
            confirmButtonColor: "#3da09a"
        })
            return
        }
        asset.reset();
        swal({
            title:"Success!",
            text: res.Message,
            type: "success",
            confirmButtonColor: "#3da09a"
        });
        $('#addNewAsset').modal('hide');
        asset.getData(function () {
            asset.renderGrid()
        })
        model.Processing(false);
    })

}

asset.editData = function () {

    var validator = $("#myForm").data("kendoValidator");
    if (validator === undefined) {
        validator = $("#myForm").kendoValidator().data("kendoValidator");
    }
    if (!validator.validate()) {
        return
    }
    var assetPrice = FormatCurrency(asset.record.Price())

    var url = "/transaction/editdataasset"

    var payload = ko.mapping.toJS(asset.record)
    payload.MonthlyDepreciation = FormatCurrency(asset.record.MonthlyDepreciation())
    payload.Price = assetPrice
    payload.Qty = FormatCurrency(asset.record.Qty())
    payload.Total = FormatCurrency(asset.record.Total())
    payload.User = userinfo.usernameh()
    payload.PostingDate = $('#postingdate').data('kendoDatePicker').value()
    payload.DatePeriod = $('#dateperiod').data('kendoDatePicker').value()

    var param = {
        // Data : JSON.stringify(payload),
        Data : payload,
        ID : asset.record.ID      
    }
    model.Processing(true)
    ajaxPost(url, param, function (res) {
        if (res.IsError === "true") {
            swal("Error!", res.Message, "error")
            return
        }
        asset.reset();
        swal({
            title:"Success!",
            text: res.Message,
            type: "success",
            confirmButtonColor:"#3da09a"});
        $('#addNewAsset').modal('hide');
        asset.getData(function () {
            asset.renderGrid()
        })
        model.Processing(false);
    })
}
asset.getDataDepreciation = function(callback){
    model.Processing(true)
    var param = {}
    ajaxPost('/transaction/getdatadepreciationasset', param, function (res) {
        model.Processing(false)
        asset.dataDepreciationMonth($.extend(true, [], res.Data))
        callback()
    }, function () {
        swal({
            title:"Error!",
            text: "Unknown error, please try again",
            type: "error",
            confirmButtonColor:"#3da09a"
        })
        model.Processing(false)
    })
}
asset.getData = function (callback) {
    model.Processing(true)
    var enddate = $('#dateEnd').data('kendoDatePicker').value();
    enddate.setMonth(enddate.getMonth()+1)
    enddate.setDate(0)
    var param = {}
    if (asset.filterindicator() == true) {
        param = {
            DateFilter: enddate,
            Filter: true,
            TextSearch: asset.textSearch(),
        }
        if (param.TextSearch == "") {
            param.Filter = false
        }
    } else {
        param = {
            Filter: false,
        }
    }
    ajaxPost('/transaction/getdataasset', param, function (res) {
        model.Processing(false)
        if (res.IsError) {
            swal("Search Not Found!", res.Message, "warning")
            $('#textSearch').val("")
            return
        }
        asset.dataMasterAsset($.extend(true, [], res.Data))
        asset.dataMasterAssetOriginal($.extend(true, [], res.Data))
        asset.dataFilter(res.Data)
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

asset.renderGrid = function () {
    var data = asset.dataMasterAsset();
    //console.log(data)
    var listcolumns = []
    var listaggregate = []

    listcolumns.push(
    {
        title: 'Action',
        width: 70,
        locked: true,
        lockable: false,
        template: "#if(1==1){#" +"<a href=\"javascript:asset.edit('#: _id #')\" data-target=\".EditAsset\" data-backdrop=\"static\" class=\"btn btn-xs btn-warning\"><i class=\"fa fa-pencil\"></i></a>&nbsp;" +
                "<a href=\"javascript:asset.delete('#: _id #', '#: Description #')\" class=\"btn btn-xs btn-danger\"><span class='glyphicon glyphicon-trash'></span>#}#",
        headerAttributes: {
            style: "vertical-align: middle;float: none;"
        }
    },
    {
        title: "No.",
        filterable: false,
        template: function (dataItem) {
            var idx = _.findIndex(data, function (d) {
                return d._id == dataItem._id
            })

            return idx + 1
        },
        width: 40,
        locked: true,
        lockable: false,
        headerAttributes: {
            style: "vertical-align: middle;float: none;"
        }
    }, {
        field: 'Description',
        title: 'Description',
        width: 200,
        locked: true,
        lockable: false,
        footerTemplate: "Total",
        headerAttributes: {
            style: "vertical-align: middle;float: none;"
        }
    }, {
        field: 'Category',
        title: 'Category',
        width: 100,
        headerAttributes: {
            style: "vertical-align: middle;float: none;"
        }
    }, {
        field: 'Qty',
        title: 'Qty',
        width: 50,
        headerAttributes: {
            style: "vertical-align: middle;float: none;"
        }

    }, {
        field: 'Price',
        title: 'Price',
        width: 200,
        template: "#=ChangeToRupiah(Price)#",
        attributes: {
            "class": "align-right"
        },
        headerAttributes: {
            style: "vertical-align: middle;float: none;"
        }
    }, {
        field: 'Total',
        title: 'Total',
        width: 200,
        template: "#=ChangeToRupiah(Total)#",
        attributes: {
            "class": "align-right"
        },
        aggregates: ["sum"],
        footerTemplate: "<div class='align-right'>#=ChangeToRupiah(sum)#</div>",
        headerAttributes: {
            style: "vertical-align: middle;float: none;"
        }

    }, {
        field: 'PostingDate',
        title: 'Posting Date',
        width: 150,
        template: "#= moment(PostingDate).format('DD-MMM-YYYY')#",
        headerAttributes: {
            style: "vertical-align: middle;float: none;"
        }
    }, {
        field: 'DatePeriod',
        title: 'Date of End Period',
        width: 150,
        template: "#= moment(DatePeriod).format('DD-MMM-YYYY')#",
        headerAttributes: {
            style: "vertical-align: middle;float: none;"
        }

    }, {
        field: 'SumDepreciation',
        title: 'Sum of months<br/> for Depreciations',
        width: 200,
        headerAttributes: {
            style: "vertical-align: middle;float: none;"
        }

    }, {
        field: 'MonthlyDepreciation',
        title: 'Monthly Depreciaton',
        width: 200,
        template: "#=ChangeToRupiah(MonthlyDepreciation)#",
        attributes: {
            "class": "align-right"
        },
        headerAttributes: {
            style: "vertical-align: middle;float: none;"
        }

    })

    var arraggtotal = {
        field: "Total",
        aggregate: "sum",
    }

    listaggregate.push(arraggtotal)
    var date = $('#dateEnd').data('kendoDatePicker').value();
    var maxYear = moment(date).year()
    var maxMonth = moment(date).month()
    if (data.length > 0) {
        var minYear = moment(_.minBy(data, function (k) {
            return k.PostingDate;
        }).PostingDate).year();
    }

    for (var i = minYear; i < maxYear; i++) {
        _.each(data, function (v, e) {
            if (data[e]['totalaccum'] === undefined) {
                data[e]['totalaccum'] = 0
            }

            if (moment(data[e].PostingDate).year() > i) {
                data[e]['accum' + i] = 0
                data[e]['totalaccum'] += data[e]['accum' + i]
                return
            }

            if (moment(data[e].PostingDate).year() == i && moment(data[e].DatePeriod).year() == i) {
                var monthvalue = moment(data[e].DatePeriod).month() - moment(data[e].PostingDate).month()
                data[e]['accum' + i] = v.Total
                data[e]['totalaccum'] += data[e]['accum' + i]
                return
            }

            if (moment(data[e].PostingDate).year() == i) {
                var monthvalue = 12 - moment(data[e].PostingDate).month()
                data[e]['accum' + i] = v.MonthlyDepreciation * monthvalue
                data[e]['totalaccum'] += data[e]['accum' + i]
                return
            }

            if (moment(data[e].DatePeriod).year() == i) {
                var monthvalue = moment(data[e].DatePeriod).month()
                data[e]['accum' + i] = v.Total - data[e]['totalaccum']
                data[e]['totalaccum'] += data[e]['accum' + i]
                return
            }

            if (i > moment(data[e].PostingDate).year() && moment(data[e].DatePeriod).year() > i) {
                data[e]['accum' + i] = v.MonthlyDepreciation * 12
                data[e]['totalaccum'] += data[e]['accum' + i]
                return
            }
            data[e]['accum' + i] = 0
            data[e]['totalaccum'] += data[e]['accum' + i]


        })

        var columne = {
            field: 'accum' + i,
            title: 'Accumulation Of <br/>Depreciation ' + i,
            width: 200,
            template: "#=ChangeToRupiah(accum" + i + ")#",
            attributes: {
                "class": "align-right"
            },
            aggregates: ["sum"],
            footerTemplate: "<div class='align-right'>#=ChangeToRupiah(sum)#</div>",
            headerAttributes: {
                style: "vertical-align: middle;float: none;"
            }
        }

        listcolumns.push(columne)

        var arragg = {
            field: 'accum' + i,
            aggregate: "sum",
        }
        listaggregate.push(arragg)
    }

    var dataIdcheckBox = []
    for (var i = 0; i <= maxMonth; i++) {
        _.each(data, function (v, e) {
            if (data[e]['summonth'] === undefined) {
                data[e]['summonth'] = 0
            }
            if (data[e]['totalaccum'] === undefined) {
                data[e]['totalaccum'] = 0
            }
            if (data[e]['totalminusaccum'] === undefined) {
                data[e]['totalminusaccum'] = data[e]['Total'] - data[e]['totalaccum']
                data[e]['sumtotalaccum'] = data[e]['totalaccum']
            }

            if (maxYear > moment(data[e].DatePeriod).year()) {
                data[e]['month' + i] = 0
                data[e]['summonth'] += data[e]['month' + i]
                return
            }

            if (maxYear < moment(data[e].DatePeriod).year()) {
                if (i >= moment(data[e].PostingDate).month()) {
                    data[e]['month' + i] = v.MonthlyDepreciation
                    data[e]['summonth'] += data[e]['month' + i]
                    return
                } else if ((i <= moment(data[e].PostingDate).month()) && (maxYear > moment(data[e].PostingDate).year())) {
                    data[e]['month' + i] = v.MonthlyDepreciation
                    data[e]['summonth'] += data[e]['month' + i]
                    return
                }
            }

            if (maxYear == moment(data[e].DatePeriod).year() && maxYear == moment(data[e].PostingDate).year()) {
                if (moment(data[e].PostingDate).month() <= i && i < moment(data[e].DatePeriod).month()) {
                    data[e]['month' + i] = v.MonthlyDepreciation
                    data[e]['totalminusaccum'] = data[e]['totalminusaccum'] - data[e]['month' + i]
                    data[e]['summonth'] += data[e]['month' + i]
                } else if (i == moment(data[e].DatePeriod).month()) {
                    data[e]['month' + i] = data[e]['totalminusaccum']
                    data[e]['summonth'] += data[e]['month' + i]
                } else {
                    data[e]['month' + i] = 0
                    data[e]['summonth'] += data[e]['month' + i]
                }
                return
            }

            if (maxYear == moment(data[e].DatePeriod).year()) {
                if (i < moment(data[e].DatePeriod).month()) {
                    data[e]['month' + i] = v.MonthlyDepreciation
                    data[e]['totalminusaccum'] = data[e]['totalminusaccum'] - data[e]['month' + i]
                    data[e]['summonth'] += data[e]['month' + i]
                } else if (i == moment(data[e].DatePeriod).month()) {
                    data[e]['month' + i] = data[e]['totalminusaccum']
                    data[e]['summonth'] += data[e]['month' + i]
                } else {
                    data[e]['month' + i] = 0
                    data[e]['summonth'] += data[e]['month' + i]
                }
                return
            }

            data[e]['month' + i] = 0
            data[e]['summonth'] += data[e]['month' + i]
        })
        var monthInt= i+1
        var dataMY = maxYear+""+monthInt
        var dataDepreciation = _.filter(asset.dataDepreciationMonth(), function(e){return e.MonthYear== dataMY})
        // console.log(asset.dataDepreciationMonth(),dataDepreciation, dataMY)
        var checkboxHeader ="<div class='align-center'><input type='checkbox' id='month"+i+"' class='checkbox' checked disabled/>"+ moment.months(i) + " " + maxYear+"</div"
        if (dataDepreciation.length == 0){
            dataIdcheckBox.push({
                "idCheckBox":'month' + i,
                "MonthYear": dataMY,
                "Title": moment.months(i) + " " + maxYear,
            })
            checkboxHeader = "<div class='align-center'><input type='checkbox' id='month"+i+"' class='checkbox'/>"+ moment.months(i) + " " + maxYear+"</div"
        }
        var columne = {
            field: 'month' + i,
            title: moment.months(i) + " " + maxYear,
            headerTemplate: checkboxHeader,
            // headerTemplate: moment.months(i) + " " + maxYear+"<input type='checkbox' id='month"+i+"' class='checkbox' #: readyInJpurnal==true ? 'checked' : ''#/>",
            width: 200,
            template: "#=ChangeToRupiah(month" + i + ")#",
            attributes: {
                "class": "align-right"
            },
            aggregates: ["sum"],
            footerTemplate: "<div class='align-right'>#=ChangeToRupiah(sum)#</div>",
            headerAttributes: {
                style: "vertical-align: middle;float: none;"
            }
        }

        listcolumns.push(columne)

        var arragg = {
            field: 'month' + i,
            aggregate: "sum",
        }
        listaggregate.push(arragg)
    }
    asset.dataIdcheckBox(dataIdcheckBox)
    _.each(data, function (v, e) {

        if (maxYear > moment(data[e].DatePeriod).year()) {
            data[e]['accum' + maxYear] = data[e]['sumtotalaccum']
            return
        }

        data[e]['accum' + maxYear] = data[e]['sumtotalaccum'] + data[e]['summonth']

    })

    var columne = {
        field: 'accum' + maxYear,
        title: 'Accumulation Of <br/>Depreciation ' + maxYear,
        width: 200,
        template: "#=ChangeToRupiah(accum" + maxYear + ")#",
        attributes: {
            "class": "align-right"
        },
        aggregates: ["sum"],
        footerTemplate: "<div class='align-right'>#=ChangeToRupiah(sum)#</div>",
        headerAttributes: {
            style: "vertical-align: middle;float: none;"
        }
    }

    listcolumns.push(columne)

    var arragg = {
        field: "accum" + maxYear,
        aggregate: "sum",
    }
    listaggregate.push(arragg)

    _.each(data, function (v, e) {
        if (moment(data[e].DatePeriod).year() >= maxYear) {
            data[e]['netbook' + maxYear] = data[e].Total - data[e]['accum' + maxYear]
            return
        }
        data[e]['netbook' + maxYear] = 0
    })

    var columne2 = {
        field: 'netbook' + maxYear,
        title: 'Net Book Value ' + maxYear,
        width: 200,
        template: "#=ChangeToRupiah(netbook" + maxYear + ")#",
        attributes: {
            "class": "align-right"
        },
        aggregates: ["sum"],
        footerTemplate: "<div class='align-right'>#=ChangeToRupiah(sum)#</div>",
        headerAttributes: {
            style: "vertical-align: middle;float: none;"
        }
    }

    listcolumns.push(columne2)

    var arragg2 = {
        field: "netbook" + maxYear,
        aggregate: "sum",
    }
    listaggregate.push(arragg2)

    $('#gridAsset').html("")
    $('#gridAsset').kendoGrid({
        dataSource: {
            data: data,
            aggregate: listaggregate,
            sort: {
                field: 'PostingDate',
                dir: 'asc',
            }
        },
        // sortable: true,
        excelExport: function (e) {
            var rows = e.workbook.sheets[0].rows;
            for (var ri = 0; ri < rows.length; ri++) {
                var row = rows[ri];
                if (row.type == "group-footer" || row.type == "footer") {
                    for (var ci = 0; ci < row.cells.length; ci++) {
                        var cell = row.cells[ci];
                        if (cell.value) {
                            // Use jQuery.fn.text to remove the HTML and get only the text
                            var text = $(cell.value).text();
                            if (cell.value != "Total") {
                                var expression = /^\([\d,\.]*\)$/;
                                if (text.match(expression)) {
                                    //It matched - strip out parentheses and append - at front
                                    var val = '-' + text.replace(/[\(\)]/g, '')
                                } else {
                                    var val = text
                                }
                                cell.value = parseFloat(val.split(",").join("")) // this different
                                cell.format = "#,##0.00_);(#,##0.00);0.00;"
                                // Set the alignment
                                cell.hAlign = "right";
                            } else {
                                cell.hAlign = "right";
                            }
                        }
                    }
                }
                if (row.type == "header") {
                    for (var ci = 0; ci < row.cells.length; ci++) {
                        var cell = row.cells[ci];
                        cell.value = cell.value.replace("<br/>", "")
                    }
                }
                if (row.type == "data") {
                    if (ri > 0) {
                        for (var ci = 0; ci < row.cells.length; ci++) {
                            var cell = row.cells[ci];
                            if (ci == 3 || ci == 4 || ci > 7) {
                                cell.format = "#,##0.00_);(#,##0.00);0.00;"
                                // Set the alignment
                                cell.hAlign = "right";
                            }

                        }
                    }
                }
            }
            e.workbook.fileName = "transaction-asset"
        },
        scrollable: true,
        columns: listcolumns,
        height: 500

    })
}

asset.search = function (e) {
    asset.textSearch(e)
    asset.filterindicator(true)
    asset.getData(function () {
        asset.renderGrid()
    })
}

asset.reset = function () {
    $("#description").val("");
    $("#category").data("kendoDropDownList").value("")
    $("#qty").val("");
    $("#price").val("");
    $("#total").val("");
    $("#postingdate").data("kendoDatePicker").value(new Date());
    $("#dateperiod").data("kendoDatePicker").value(new Date());
    $("#sum").val("");
    $("#depreciation").val("");
}

asset.getDataCategory = function () {
    var param = {}
    ajaxPost('/master/getdatacategory', param, function (res) {
        if (res.IsError === "true") {
            swal({
                title:"Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor:"#3da09a"})
            return
        }
        asset.dataMasterCategory(res.Data)
    })
}

asset.maskingMoney = function () {
    $('.currency').inputmask("numeric", {
        radixPoint: ".",
        groupSeparator: ",",
        digits: 2,
        autoGroup: true,
        rightAlign: false,
    });
}
asset.ImportExcel = function () {
    event.stopPropagation();
    event.preventDefault();

    model.Processing(true);

    if ($('input[type=file]')[0].files[0] == undefined) {
        swal('Error', 'Please select a file to upload!', 'error');
        model.Processing(false);
        return;
    }
    var len = $("#fDok")[0].files.length
    if (len > 0) {
        var j = 0;
        var data = new FormData();
        for (i = 0; i < $("#fDok")[0].files.length; i++) {
            data.append("filedoc", $('input[type=file]')[0].files[i]);
            data.append("filename", $('input[type=file]')[0].files[i].name);
        }
    } else {
        var data = new FormData();
        data.append("filedoc", $('input[type=file]')[0].files[0]);
        data.append("filename", $('input[type=file]')[0].files[0].name);
    }
    if ($('input[type=file]')[0].files[0].name != "") {
        jQuery.ajax({
            url: '/transaction/importexcelasset',
            data: data,
            cache: false,
            contentType: false,
            processData: false,
            type: 'POST',
            success: function (data) {
                swal('Success', 'Data has been uploaded successfully!', 'success');
                model.Processing(false);
                $('#fDok').val('');
                $('#fInfo').val('');
                $('#ImportModal').modal('hide');
                asset.getData(function () {
                    asset.renderGrid()
                })
            }
        });
    } else {
        swal('Error', 'Please select a file to upload!', 'error');
        model.Processing(false);
    }
}

asset.ValidateFiles = function (elm) {
    var file = elm.val();
    var ext = file.split(".");
    ext = ext[ext.length - 1].toLowerCase();
    var arrayExtensions = ["xls", "xlsx", "csv"];

    if (arrayExtensions.lastIndexOf(ext) == -1) {
        swal('Error', 'Invalid file extension. Please upload an excel file format!', 'error');
        elm.val('');
    }
}

asset.exportExcel = function () {
    $("#gridAsset").getKendoGrid().saveAsExcel();
}

asset.saveEdit = function(){
    if (asset.button() == "Update"){
        //console.log("DIMARI2")
        asset.editData()
        return
    }
    asset.saveData()
}

asset.edit = function(e){
    //console.log(e)
    var datareal = asset.dataFilter()    
    var data = _.find(datareal, function (item) {
        return item._id == e;
    });
    ko.mapping.fromJS(data, asset.record)
    asset.record.Total = ko.computed(function () {
        var n = String(asset.record.Price())
        var total = asset.record.Qty() * parseFloat(n.replace(/,/g, ''))
        return total.toFixed(2);
    })
    asset.record.SumDepreciation = ko.computed(function () {
        var month = 0
        var a = moment(asset.record.PostingDate())
        var b = moment(asset.record.DatePeriod())
        month = b.diff(a, 'month')
        return (month + 1);
    })

    asset.record.MonthlyDepreciation = ko.computed(function () {
        var month = 0
        var a = moment(asset.record.PostingDate())
        var b = moment(asset.record.DatePeriod())
        month = b.diff(a, 'month')
        var total = asset.record.Total() / (month + 1)
        return total.toFixed(2);
    })
    //console.log(asset.record)
    asset.record.ID(e)
    $("#addNewAsset").modal("show")  
    asset.button("Update")  
    asset.textTitle("Edit Asset")

}

asset.delete = function(e,desc){
    var url = "/transaction/deletedataasset"
    var param = {
        ID: e
    }
    swal({
        title: "Are you sure?",
        text: "You will delete " + desc ,
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3da09a",
        confirmButtonText: "Yes, do it!",
        cancelButtonText: "No!",
        closeOnConfirm: true,
        closeOnCancel: false
    }, function (isConfirm) {
        if (isConfirm) {
            model.Processing(true)
            ajaxPost(url, param, function (res) {
                if (res.IsError == true) {
                    model.Processing(false)
                    return swal("Error!", "Error to delete this Asset!!!", "error")
                }
                setTimeout(function () {
                    asset.reset();
                    swal("Success!", res.Message, "success");
                    $('#addNewAsset').modal('hide');
                    asset.getData(function () {
                        asset.renderGrid()
                    })
                    model.Processing(false);
                }, 100);

            })
        } else {
            swal({
                title:"Cancelled",
                type:"error",
                confirmButtonColor:"#3da09a"});
        }
    });
}
asset.buttonProcessDepreciation= function(){
    var dataModal = []
    var dataCheckBox = asset.dataIdcheckBox()
    _.each(dataCheckBox, function(v,i){
        var checked = $('#'+v.idCheckBox+':checkbox:checked').length > 0;
        if (checked){
            var sumPerDep = $("#gridAsset").data("kendoGrid").dataSource.aggregates()
            var amount = sumPerDep[v.idCheckBox].sum
            dataModal.push({
                _id:"",
                Amount: amount,
                Date: new Date(),
                DateMonthYear: v.Title,
                IdChecbox: v.idCheckBox,
                MonthYear: v.MonthYear,
                Checked: true
            })
        }
    })
    // console.log(dataModal)
    asset.dataModalDepreciation(dataModal)
    $('#processAssetJournal').modal('show');
    // var atLeastOneIsChecked = $('#checkArray:checkbox:checked').length > 0;
}
asset.onChangeDate= function(value, index, monthyear){
    // console.log(value, index, monthyear)
    var data = asset.dataModalDepreciation()
    _.each(data, function(v,i){
        if(monthyear==v.MonthYear){
            data[i].Date= value
        }
    })
    asset.dataModalDepreciation(data)
}
asset.saveDataDepreciation = function(){
    var url = "/transaction/savedatadepreciationandjournal"
    var param = {
        Data : asset.dataModalDepreciation()
    }
    swal({
        title: "Are you sure?",
        text: "Save this data for journal",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3da09a",
        confirmButtonText: "Yes, do it!",
        cancelButtonText: "No!",
        closeOnConfirm: true,
        closeOnCancel: false
    }, function (isConfirm) {
        if (isConfirm) {
            model.Processing(true)
            ajaxPost(url, param, function (res) {
                // console.log(res)
                $('#processAssetJournal').modal('hide');
                if (res.IsError === "true") {
                    swal({
                        title:"Error!",
                        text: res.Message,
                        type: "error",
                        confirmButtonColor: "#3da09a"
                    })
                    return
                }
                setTimeout(function(){
                    swal({
                        title: "Success!",
                        text: "Data has been saved!",
                        type: "success",
                        confirmButtonColor: "#3da09a",
                    }, function () {
                        model.Processing(false);
                        location.reload()
                    });
                },100)
                model.Processing(false);
            })
        } else {
            swal({
                title:"Cancelled",
                type:"error",
                confirmButtonColor:"#3da09a"
            });
        }
    });
    // ajaxPost(url, param, function(res){
    //     if (res.IsError === "true") {
    //         swal({
    //         title:"Error!",
    //         text: res.Message,
    //         type: "error",
    //         confirmButtonColor: "#3da09a"
    //     })
    //         return
    //     }
    //     asset.reset();
    //     swal({
    //         title: "Success!",
    //         text: "Data has been saved!",
    //         type: "success",
    //         confirmButtonColor: "#3da09a",
    //     }, function () {
    //         model.Processing(false);
    //         // location.reload()
    //     });

    //     $('#processAssetJournal').modal('hide');
    //     model.Processing(false);
    // })
    // console.log(asset.dataModalDepreciation())
}
asset.init = function () {
    asset.getDateNow()
    asset.maskingMoney()
    asset.getDataCategory()
    asset.getData(function () {
        asset.getDataDepreciation(function(){
            asset.renderGrid()
        })
    })
}

$(function () {
    $(document).on('change', ':file', function () {
        asset.ValidateFiles($(this));
        var input = $(this),
            numFiles = input.get(0).files ? input.get(0).files.length : 1,
            label = input.val().replace(/\\/g, '/').replace(/.*\//, '');
        input.trigger('fileselect', [
        numFiles, label]);
    });
    $(':file').on('fileselect', function (event, numFiles, label) {
        var input = $(this).parents('.input-group').find(':text'),
            log = numFiles > 1 ? numFiles + ' files selected' : label;
        if (input.length) {
            input.val(log);
        } else {
            if (log) alert(log);
        }
    });
    asset.init()
})

