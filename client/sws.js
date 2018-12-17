var d = document
var w = window
var l = window.location
var n = w.navigator
var esc = encodeURIComponent

function send (p, obj) {
  console.log('sending', p, JSON.stringify(obj))
  var qs = Object.keys(obj)
    .map(function(k) {
      return esc(k) + '=' + esc(obj[k])
    })
    .join('&')
  var r = new w.XMLHttpRequest()
  r.open('GET', p + "?" + qs, true)
  r.send()
}

/*
function ch (e) {
  var href
  var el = e.target || e.srcElement
  if (el.tagName === 'A') {
    href = el.getAttribute('href')
  }
  var rect = el.getBoundingClientRect()
  send('/click', {
    url: d.URL,
    href: href,
    top: rect.top,
    right: rect.right,
    bottom: rect.bottom,
    left: rect.left
  })
}

function frameBuster () {
  if (w.location !== w.top.location) {
    w.top.location = w.location
  }
}
*/

function ready (fn) {
  if (d.attachEvent
    ? d.readyState === 'complete'
    : d.readyState !== 'loading') {
    fn()
    /*
    if (_sws.events) {
      d.attachEvent('onclick', ch)
    }
    */
  } else {
    d.addEventListener('DOMContentLoaded', fn)
    /*
    if (_sws.events) {
      d.addEventListener('click', ch)
    }
    */
  }
}

/*
function detectBrowserFeatures() {
  var i,
    mimeType,
    pluginMap = {
      // document types
      pdf: 'application/pdf',

      // media players
      qt: 'video/quicktime',
      realp: 'audio/x-pn-realaudio-plugin',
      wma: 'application/x-mplayer2',

      // interactive multimedia
      dir: 'application/x-director',
      fla: 'application/x-shockwave-flash',

      // RIA
      java: 'application/x-java-vm',
      gears: 'application/x-googlegears',
      ag: 'application/x-silverlight'
    };

  // detect browser features except IE < 11 (IE 11 user agent is no longer MSIE)
  if (!((new RegExp('MSIE')).test(n.userAgent))) {
    // general plugin detection
    if (n.mimeTypes && n.mimeTypes.length) {
      for (i in pluginMap) {
        if (Object.prototype.hasOwnProperty.call(pluginMap, i)) {
          mimeType = n.mimeTypes[pluginMap[i]];
          browserFeatures[i] = (mimeType && mimeType.enabledPlugin) ? '1' : '0';
        }
      }
    }
    // Safari and Opera
    // IE6/IE7 navigator.javaEnabled can't be aliased, so test directly
    // on Edge navigator.javaEnabled() always returns `true`, so ignore it
    if (!((new RegExp('Edge[ /](\\d+[\\.\\d]+)')).test(n.userAgent)) &&
      typeof navigator.javaEnabled !== 'unknown' &&
      isDefined(n.javaEnabled) &&
      n.javaEnabled()) {
      browserFeatures.java = '1';
    }

// Firefox
    if (isFunction(windowAlias.GearsFactory)) {
      browserFeatures.gears = '1';
    }

// other browser features
    browserFeatures.cookie = hasCookies();
  }
}

var width = parseInt(screenAlias.width, 10);
var height = parseInt(screenAlias.height, 10);
browserFeatures.res = parseInt(width, 10) + 'x' + parseInt(height, 10);
*/

var viewPort = (w.innerWidth || d.documentElement.clientWidth || d.body.clientWidth)
viewPort += "x"
viewPort += (w.innerHeight || d.documentElement.clientHeight || d.body.clientHeight)

ready(function () {
  send('http://localhost:5000/sws.gif', {
    s: l.protocol,
    h: l.host,
    p: l.pathname + l.search + l.hash, // page
    t: _sws.title || d.title,
    r: d.referrer,
    u: n.userAgent,
    v: viewPort,
  })
})
