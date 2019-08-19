/*
 * @Author: Ainur
 * @Date:   2016-11-17 13:40:31
 * @Last Modified by:   Ainur
 * @Last Modified time: 2016-11-23 07:34:29
 */

 function debounce(func, wait, immediate) {
  var timeout;
  return function() {
    var context = this, args = arguments;
    var later = function() {
      timeout = null;
      if (!immediate) func.apply(context, args);
    };
    var callNow = immediate && !timeout;
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
    if (callNow) func.apply(context, args);
  };
};

var search = kendo.observable({
  SelectedItem: "",
  Processing: false,
  Data: [],
  keypressed: function(e){
    if (e.which == 13) {
      search.GetData(e.target.value)
    } 
  }
})

search.GetData = function (keyword) {
//  console.log("SEARCH");
  var url = "/search/getresult";
  var parm = {};
  if (keyword === undefined) {
    parm["keyword"] = search.get("SelectedItem").Keyword;
  } else {
    parm["keyword"] = keyword;
  }

  parm["keywordEscaped"] = parm["keyword"]
    .replace(/\(/g, "\\(")
    .replace(/\)/g, "\\)")
    .replace(/\+/g, "\\+")
    .replace(/\-/g, "\\-")
    .replace(/\//g, "\\/")
    .replace(/\./g, "\\.")
    .replace(/\&/g, "\\&")
    .replace(/\,/g, "\\,")

  ajaxPost(url, parm, function (res) {
    if (res.Data.ContractName.length === 0 && res.Data.ContractManager.length === 0 &&
            res.Data.GoverningLaw.length === 0) {
      swal("Error!", "No Matching on choosen keyword.", "error");
      return false;
    } else {
      $("#search-result").modal("show");
      search.Render(res.Data);
    }
  });
}

search.Render = function (dataSource) {
//  console.log(dataSource);
  $("#sr-data").html("");
  $("#sr-data").kendoPanelBar({
    expandMode: "multiple"
  });
  var dataArr = [];
  var pnList = [];
  var pmList = [];
  var aeList = [];

  if (dataSource.ContractName.length > 0) {
    $.each(dataSource.ContractName, function (idx, obj) {
      pnList.push({
        "encoded": false, 
        "text": "<a href='/dashboard/review?id=" + obj.Id + "' style='text-decoration: underline;'> " + obj.Keyword + "</a>", 
        "value": obj.Id
      });
    });

    dataArr.push({
      text: "<b>Contract Name ( " + dataSource.ContractName.length + " )</b>", 
      value: "VOID", 
      encoded: false, 
      items: pnList
    });
  }

  if (dataSource.ContractManager.length > 0) {
    $.each(dataSource.ContractManager, function (idx, obj) {
      pmList.push({
        "encoded": false, 
        "text": "<table style='width:100%'> <tbody><tr><td style='width:50%; vertical-align:top;'><a href='/dashboard/review?id=" + obj.Id + "' style='text-decoration: underline;'>" + obj.Keyword + "</a></td><td style='width:50%'>" + obj.ParentInfo + "</td></tr></tbody></table>", 
        "value": obj.Id
      });
    });

    dataArr.push({text: "<b>Contract Manager ( " + dataSource.ContractManager.length + " )</b>", value: "VOID", encoded: false, items: pmList});
  }

  if (dataSource.GoverningLaw.length > 0) {
    $.each(dataSource.GoverningLaw, function (idx, obj) {
      aeList.push({
        "encoded": false, 
        "text": "<table style='width:100%'> <tbody><tr><td style='width:50%; vertical-align:top;'><a href='/dashboard/review?id=" + obj.Id + "' style='text-decoration: underline;'>" + obj.Keyword + "</a></td><td style='width:50%'>" + obj.ParentInfo + "</td></tr></tbody></table>", 
        "value": obj.Id
      });
    });

    dataArr.push({text: "<b>Receiving Entity ( " + dataSource.GoverningLaw.length + " )</b>", value: "VOID", encoded: false, items: aeList});
  }

  var srdata = $("#sr-data").data("kendoPanelBar");
  srdata.append(dataArr);
};

$(document).ready(function (){
  kendo.bind($("#bredkram"), search)
});