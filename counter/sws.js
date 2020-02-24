var d = document
var de = d.documentElement
var w = window
var l = d.location
var n = w.navigator
var esc = encodeURIComponent
var me = d.currentScript

var _sws = w._sws || {noauto: false, local: false}
_sws.d = _sws.d || me.src
_sws.site = _sws.site || me.dataset.site

function count (p, obj) {
  if (!_sws.local && l.hostname.match(/(localhost$|^127\.|^10\.|^172\.16\.|^192\.168\.)/))
    return
  if ('visibilityState' in d && d.visibilityState === 'prerender')
    return

  var qs = Object.keys(obj)
    .map(function (k) {
      return esc(k) + '=' + esc(obj[k])
    })
    .join('&')

  if (!_sws.noxhr) {
    var r = new w.XMLHttpRequest()
    r.open('GET', p + '?' + qs, true)
    r.send()
  } else {
    var image = new Image(1, 1)
    image.src = p + '?' + qs
  }
}

function ready (fn) {
  if (d.attachEvent ? d.readyState === 'complete' : d.readyState !== 'loading') {
    fn()
  } else {
    d.addEventListener('DOMContentLoaded', fn)
  }
}

var viewPort = (w.innerWidth || de.clientWidth || d.body.clientWidth) + 'x' +
  (w.innerHeight || de.clientHeight || d.body.clientHeight)

ready(function () {
  if (!_sws.noauto) {
    var ep = new URL(_sws.d)
    count('{{ .endpoint }}', {
      i: _sws.site,
      s: l.protocol,
      h: l.host,
      p: l.pathname,
      q: l.search + l.hash,
      t: _sws.title || d.title,
      r: d.referrer,
      u: n.userAgent,
      v: viewPort
    })
  }
})
