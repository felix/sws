var d = document
var de = document.documentElement
var w = window
var l = window.location
var n = w.navigator
var esc = encodeURIComponent
var me = document.currentScript
console.log('me:', me)
console.log('me.sws:', me.dataset.sws)

_sws = _sws || {xhr:true}
console.log('_sws:', _sws)
_sws.d = _sws.d || me.dataset.sws || 'http://sws.userspace.com.au/sws.gif'
console.log('using', _sws.d)

function send (p, obj) {
  console.log('sending', p, JSON.stringify(obj))
  var qs = Object.keys(obj)
    .map(function (k) {
      return esc(k) + '=' + esc(obj[k])
    })
    .join('&')
  if (_sws.xhr) {
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

var viewPort = (w.innerWidth || de.clientWidth || d.body.clientWidth)
  + 'x'
  + (w.innerHeight || de.clientHeight || d.body.clientHeight)

ready(function () {
  send(_sws.d, {
    s: l.protocol,
    h: l.host,
    p: l.pathname,
    q: l.search + l.hash,
    t: _sws.title || d.title,
    r: d.referrer,
    u: n.userAgent,
    v: viewPort
  })
})
