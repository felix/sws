var d = document
var de = document.documentElement
var w = window
var l = d.location
var n = w.navigator
var esc = encodeURIComponent
var me = document.currentScript
console.log('me:', me)
console.log('me.sws:', me.dataset.sws)

var _sws = w._sws || {noxhr: false, noauto: false}
console.log('_sws:', _sws)
_sws.d = _sws.d || me.dataset.sws || 'http://sws.userspace.com.au/sws.gif'
_sws.site = _sws.site || me.dataset.site
console.log('using', _sws.d)

function count (p, obj) {
  console.log('sending', p, JSON.stringify(obj))
  var qs = Object.keys(obj)
    .map(function (k) {
      return esc(k) + '=' + esc(obj[k])
    })
    .join('&')
  if (!_sws.xhr) {
    var r = new w.XMLHttpRequest()
    r.open('GET', p + '?' + qs, true)
    r.send()
  } else {
    var image = new Image(1, 1)
    image.src = p + '?' + qs
  }
}

function ready (fn) {
  if (d.attachEvent
    ? d.readyState === 'complete'
    : d.readyState !== 'loading') {
    fn()
  } else {
    d.addEventListener('DOMContentLoaded', fn)
  }
}

var viewPort = (w.innerWidth || de.clientWidth || d.body.clientWidth) + 'x' +
  (w.innerHeight || de.clientHeight || d.body.clientHeight)

ready(function () {
  if (!_sws.noauto) {
    count(_sws.d, {
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
