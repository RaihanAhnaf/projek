// == ProActive global js ==
var ProActive = ProActive || {};

/* ======= Configuration ======= */
ProActive.MaxUploadSize = 1024 * 1024 * 1; // 1 MiB


/* ======= Functions ======= */
ProActive.checkUploadFilesize = function (size, swal) {
    var r = size <= this.MaxUploadSize;
    if (r) return true;
    if (swal !== undefined) {
        if (typeof swal !== "string")
            this.swalWarning("Attachment file size cannot exceed " + this.getMaxUploadFileSize() + "!")
        else
            this.swalWarning(swal.replace("$size", this.getMaxUploadFileSize()));
    }
    return false;
}

ProActive.getMaxUploadFileSize = function () {
    var sz = ["bytes", "KiB", "MiB", "GiB", "TiB", "PiB"];
    var szi = 0;
    var size = this.MaxUploadSize;
    while (size >= 1024) {
        size /= 1024;
        szi++;
    }
    return (Math.round(size * 100) / 100) + " " + sz[szi];
}

ProActive.url = function (method, action, getData) {
    if (typeof action == "object") getData = action;
    if (typeof action != "string") action = "";
    if (typeof method != "string") method = "";
    if (method.length == 0 || method.substr(method.length - 1) != "/") method += "/";

    var getQry = "";
    for (var k in getData) {
        getQry += (getQry == "" ? "?" : "&") + k + "=" + getData;
    }

    return method.toLowerCase() + action.toLowerCase() + getQry;
}

ProActive.kendoExcelRender = function (e, name, rowCallback) {
    for (var si = 0; si < e.workbook.sheets.length; si++) {
        var rows = e.workbook.sheets[si].rows;

        for (var ri = 0; ri < rows.length; ri++) {
            var row = rows[ri];

            if (row.type == "header") {
                for (var ci = 0; ci < row.cells.length; ci++) {
                    var cell = row.cells[ci];
                    cell.value = cell.value.replace("<br/>", " ").replace("\t", " ");
                }
            }

            if (typeof rowCallback == "function") rowCallback(row, e.workbook.sheets[si]);
        }
    }

    if (name)
        e.workbook.fileName = name;
}

/* ======= Dialogs ======= */
ProActive.swalWarning = function (message, title) {
    return swal({
        title: title || "Warning!",
        text: typeof message == "string" ? message : "",
        type: "info",
        confirmButtonColor: "#3da09a"
    });
}

ProActive.swalConfirm = function (message, title, callback) {
    return swal({
        title: title || "Are you sure?",
        text: typeof message == "string" ? message : "",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3da09a",
        confirmButtonText: "Yes, do it!",
        cancelButtonText: "No!",
        closeOnConfirm: true,
        closeOnCancel: false
    }, callback);
}

ProActive.swalConfirmSubmit = function (action, cbConfirm, cbCancel) {
    return this.swalConfirm("You will submit this " + action, null, function (isConfirm) {
        if (isConfirm) {
            if (typeof cbConfirm == "function")
                cbConfirm();
        }
        else {
            if (typeof cbCancel == "function")
                cbCancel();
            else ProActive.swalError(null, "Cancelled");
        }
    });
}

ProActive.swalConfirmDelete = function (data, callbackYes) {
    return this.swalConfirm("Your will not be able to recover this data", "Are you sure to delete " + data + "?", function (isConfirm) {
        if (isConfirm) {
            if (typeof callbackYes == "function")
                callbackYes();
        }
    }.bind(this));
}

ProActive.swalError = function (message, title) {
    return swal({
        title: title || "Error",
        text: typeof message == "string" ? message : "",
        type: "error",
        confirmButtonColor: "#3da09a"
    });
}

ProActive.cloneData = function (source, destination) {
    if (destination === undefined || destination === null)
        destination = {};
    if (typeof source != "object" || typeof destination != "object") return false;
    for (var key in source) {
        if (typeof source[key] != "object" && typeof source[key] != "undefined" && typeof source[key] != "function") {
            destination[key] = source[key];
        } else if (typeof source[key] == "object") {
            this.cloneData(source[key], destination[key]);
        }
    }
    return destination;
}

ProActive.assignData = function (source, destination, onlyExists) {
    if (typeof source != "object" || typeof destination != "object") return false;
    onlyExists = onlyExists === true || onlyExists === undefined;
    for (var key in source) {
        if (typeof source[key] != "object" && typeof source[key] != "undefined" && typeof source[key] != "function" && (!onlyExists || key in destination)) {
            destination[key] = source[key];
        } else if (typeof source[key] == "object") {
            this.assignData(source[key], destination[key]);
        }
    }
    return true;
}

/* ======= UI Related ======= */
function ViewModel() {
    this.initialize.apply(this, arguments);
}

ViewModel.prototype.initialize = function (id) {
    this._id = id;
    this.RecordModel = {};
    this.GridID = "dataGrid";
    this.GridColumns = [];
    this.GridDataSource = {};
}

ViewModel.prototype.createRecord = function (data) {
    var model = ProActive.cloneData(this.RecordModel);
    if (data) ProActive.assignData(data, model);
    return model;
}

ViewModel.prototype.switchToList = function () { }
ViewModel.prototype.viewList = function () {
    this.renderGrid();
    this.switchToList();
}

ViewModel.prototype.switchToEditor = function () { }
ViewModel.prototype.viewRead = function () {
    this.switchToEditor();
}
ViewModel.prototype.viewEdit = function () {
    this.switchToEditor();
}

ViewModel.prototype.beforeRenderGrid = function () {
    return true;
}

ViewModel.prototype.renderGrid = function (kendoGridOptions) {
    if (!this.beforeRenderGrid()) return;
    $("#" + this.GridID).html("");
    kendoGridOptions = kendoGridOptions || {};
    var opt = {
        dataSource: {
            data: report.dataGridreport(),
            sort: {
                field: 'Date',
                dir: 'asc',
            }
        },
        excel: {
            fileName: this._id + ".xlsx"
        },
        height: 500,
        width: 140,
        scrollable: true,
        columns: this.GridColumns
    };

    $("#" + this.GridID).kendoGrid(opt)
}

ViewModel.prototype.exportGridToExcel = function (fileName) {
    var kGrid = $("#" + this.GridID).getKendoGrid();
    if (kGrid) {
        if (typeof fileName == "string" && kGrid.options && kGrid.excel)
            kGrid.options.excel.fileName = fileName + ".xlsx";
        if (kGrid.saveAsExcel)
            kGrid.saveAsExcel();
        else
            console.error("#" + this.GridID + " is not a Kendo Grid!");
    }
}

ViewModel.prototype.beforeShow = function () { }
ViewModel.prototype.afterShow = function () { }
ViewModel.prototype.show = function () {
    this.beforeShow();
    this.switchToList();
    this.afterShow();
}

ViewModel.prototype.init = function () { }

// start - the ID of KendoDatePicker start date
// end - the ID of KendoDatePicker end date
// behavior - You can set this to "end", "start" or "both" to change the behavior
ProActive.KendoDatePickerRange = function (start, end, behavior) {
    start = start || "dateStart";
    end = end || "dateEnd";
    behavior = behavior == "both" ? "both" : behavior == "start" ? "start" : "end";
    var dtpStart = $('#' + start).data('kendoDatePicker');
    var dtpEnd = $('#' + end).data('kendoDatePicker');
    dtpStart.value(moment().startOf('month').toDate());
    dtpEnd.value(moment().startOf('day').toDate());
    dtpStart._rangeBehavior = behavior;

    var refFunc = function RefreshDatePicker(dtpStart, dtpEnd) {
        var behavior = dtpStart._rangeBehavior;
        behavior = behavior == "both" ? "both" : behavior == "start" ? "start" : "end";
        if (behavior == "start" || behavior == "both") {
            dtpStart.max(dtpEnd.value());
            if (dtpStart.value() > dtpEnd.value())
                dtpStart.value(dtpEnd.value());
        }
        else
            dtpStart.max(new Date(8640000000000000));
        if (behavior == "end" || behavior == "both") {
            dtpEnd.min(dtpStart.value());
            if (dtpEnd.value() < dtpStart.value())
                dtpEnd.value(dtpStart.value());
        }
        else
            dtpEnd.min(new Date(-8640000000000000));
    }.bind(this, dtpStart, dtpEnd);

    refFunc();
    if (dtpStart._refreshRange === undefined && dtpEnd._refreshRange === undefined) {
        dtpStart._refreshRange = refFunc;
        dtpStart.bind("change", refFunc);
        dtpEnd._refreshRange = refFunc;
        dtpEnd.bind("change", refFunc);
    }
}

ProActive.KendoDatePickerChangeRangeBehavior = function(start, behavior) {
    if (start != undefined && behavior == undefined) {
        behavior = start == "both" ? "both" : start == "start" ? "start" : "end";
        start = "dateStart";
    } else {
        start = start || "dateStart";
        behavior = behavior == "both" ? "both" : behavior == "start" ? "start" : "end";
    }
    var dtpStart = $('#' + start).data('kendoDatePicker');
    dtpStart._rangeBehavior = behavior;
    if (typeof dtpStart._refreshRange == "function") {
        dtpStart._refreshRange();
    }
}

ProActive.GlobalSearch = function(gridID, fields, triggerKey) {
    triggerKey = triggerKey || "Enter";
    if (typeof gridID == "string" && typeof fields == "object") {
        var gs = {};
        gs.fields = fields;
        gs.triggerKey = triggerKey;
        gs.gridID = gridID;
        gs.globalFilter = function() {
            var filter = { logic: "or", filters: [] };
            var filteredFields = this.fields;
            $searchValue = $("#textSearch").val();
            if ($searchValue) {
                for (var k in filteredFields)
                {
                    var spl = filteredFields[k].split(" ", 2);
                    if (spl.length == 1)
                        spl[1] = "contains";
                    else
                        spl[1] = spl[1].trim();
                    spl[0] = spl[0].trim();
                    filter.filters.push({ field: spl[0], operator: spl[1], value:$searchValue});
                }
            } 
            $("#" + this.gridID).data("kendoGrid").dataSource.query({ filter: filter });
        };

        var _pool = $("#textSearch").data("GlobalSearchPool") || [];
        if (_pool.length == 0)
        {
            $("#textSearch").on("keyup blur change", function (e) {
                var pool = $("#textSearch").data("GlobalSearchPool") || [];
                for(var k in pool) {
                    var gso = pool[k];
                    if ((e.type == "keyup" && (e.key == gso.triggerKey || gso.triggerKey == "any")) || e.type != "keyup")
                    {
                        if (model && typeof model.processing == "function") {
                            model.Processing(true);
                        }
                        if (typeof gso.globalFilter == "function") {
                            gso.globalFilter();
                        }
                        if (model && typeof model.processing == "function") {
                            model.Processing(false);
                        }
                    }
                }
            });
        }
        _pool.push(gs);
        $("#textSearch").data("GlobalSearchPool", _pool);
    }
}