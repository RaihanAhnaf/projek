var captureSelection = function() {
    var text = "";
    if (document.getElementById("pdfframe").contentWindow.getSelection) {
        text = document.getElementById("pdfframe").contentWindow.getSelection().toString();
    } else if (document.selection && document.selection.type != "Control") {
        text = document.selection.createRange().text;
    }
    selectedText = text;
    
    return (document.getElementById("pdfframe").contentWindow.getSelection) ? 
    document.getElementById("pdfframe").contentWindow.getSelection().toString() : (
            (document.selection && document.selection.type != "Control") ? 
            document.selection.createRange().text : ""
        )
}

var unSelect = function(){
    if (document.getElementById("pdfframe").contentWindow.getSelection) {
        if (document.getElementById("pdfframe").contentWindow.getSelection().empty) {
            document.getElementById("pdfframe").contentWindow.getSelection().empty();
        } else if (document.getElementById("pdfframe").contentWindow.getSelection().removeAllRanges) {
            document.getElementById("pdfframe").contentWindow.getSelection().removeAllRanges();
        }
    } else if (document.getElementById("pdfframe").contentWindow.document.selection) {
        document.getElementById("pdfframe").contentWindow.document.selection.empty();
    }
}

$('#pdfframe').load(function(){
    _.each(document.getElementById("pdfframe").contentWindow.document.body.getElementsByTagName("img"), function(selector){
        $(selector).hide()
    });

    var getSelectedText = function() {
        var text = "";
        if (typeof document.getElementById("pdfframe").contentWindow.getSelection != "undefined") {
            text = document.getElementById("pdfframe").contentWindow.getSelection().toString();
        } else if (typeof document.getElementById("pdfframe").contentWindow.document.selection != "undefined" 
            && document.getElementById("pdfframe").contentWindow.document.selection.type == "Text") {
            text = document.getElementById("pdfframe").contentWindow.document.selection.createRange().text;
        }
        return text;
    }

    var doSomethingWithSelectedText = function() {
        var selectedText = getSelectedText();
        if (selectedText) {
            if(viewerVM.get("highlightClause") != ""){
                modalReview.set("currentClause", modalTraining.get("currentClause"))
                $("#modal-review").modal('show')
            }
        }
    }

    document.getElementById("pdfframe").contentWindow.document.onmouseup = doSomethingWithSelectedText;
});