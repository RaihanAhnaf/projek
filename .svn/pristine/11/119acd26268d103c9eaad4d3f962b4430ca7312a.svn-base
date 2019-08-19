var EACIIT_ATTR_ID = 'data-eaciit-id'
var EACIIT_ATTR_START = 'data-eaciit-start'
var EACIIT_ATTR_END = 'data-eaciit-end'
var HIGHLIGHT_COLOR = 'blue'

// Utils

function loadJquery() {
  var jq = document.getElementById("pdfframe").contentWindow.document.createElement('script');
  jq.src = "https://ajax.googleapis.com/ajax/libs/jquery/2.1.4/jquery.min.js";
  document.getElementById("pdfframe").contentWindow.document.getElementsByTagName('head')[0].appendChild(jq);
}

function findByEaciitId(id) {
  return document.querySelector('[' + EACIIT_ATTR_ID + '="' + id + '"]')
}

function isText(el) {
  return el.nodeType === 3
}

function getEaciitId(element) {
  var withID = isText(element) ? element.parentElement : element
  var eaciitID = withID.attributes.getNamedItem(EACIIT_ATTR_ID)
  return eaciitID === null ? null : eaciitID.value
}

function rangeForElement(element) {
  return [element.attributes[EACIIT_ATTR_START].value, element.attributes[EACIIT_ATTR_END].value]
}

function doGetSelection() {
  var range = document.getElementById("pdfframe").contentWindow.getSelection().getRangeAt(0)
  
  return {
    startElementId: getEaciitId(range.startContainer),
    startOffset: range.startOffset,
    endElementId: getEaciitId(range.endContainer),
    endOffset: range.endOffset,
    commonAncestorId: getEaciitId(range.commonAncestorContainer)
  }
}

function crawlChildren (element, fn, context) {
  var children = element.childNodes
  for (var i = 0; i < children.length; i++) {
    fn(children[i], context)
  }
}

function highlightClause(clause) {
  var commonAncestor = findByEaciitId(clause.commonAncestorId)
  var context = new Context(clause.startElementId, clause.startOffset,
    clause.endElementId, clause.endOffset, clause.commonAncestorId)
  crawlForClauseNormal(commonAncestor, context, false)
  return context
}

function Context(startElementId, startOffset, endElementId, endOffset, commonAncestorId) {
  this.startElementId = startElementId
  this.startOffset = startOffset
  this.endElementId = endElementId
  this.endOffset = endOffset
  this.commonAncestorId = commonAncestorId
  this.inRange = false
  this.words = []
}

Context.prototype.addWords = function (words) {
  this.words = this.words.concat(words)
}

function crawlForClause(element, context, lastDance) {
  if (isText(element) && (context.inRange || lastDance)) {
    var trimmed = $.trim(element.nodeValue)
    if (trimmed === '') return
    var parentRange = rangeForElement(element.parentElement)
    var startRange = 0
    var endRange = parentRange[1] - parentRange[0]
    if (getEaciitId(element.parentElement) == context.startElementId) {
      startRange = context.startOffset
    }
    if (getEaciitId(element.parentElement) == context.endElementId) {
      endRange = endRange - context.endOffset
    }
    context.addWords($.trim(element.nodeValue).split(/(\s+)/).filter(w => $.trim(w) !== ''))
    element.parentElement.style.color = HIGHLIGHT_COLOR
  } else {
    if (!context.inRange) {
      if (getEaciitId(element) == context.startElementId) {
        context.inRange = true
      }
    }
    if (context.inRange) {
      if (getEaciitId(element) == context.endElementId) {
        context.inRange = false
      }
      crawlChildren(element, crawlForClauseLastDance, context)
      return
    }
    crawlChildren(element, crawlForClauseNormal, context)
  }
}

function crawlForClauseNormal(element, context) {
  return crawlForClause(element, context, false)
}

function crawlForClauseLastDance(element, context) {
  return crawlForClause(element, context, true)
}
