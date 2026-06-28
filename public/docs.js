var btn = document.getElementById('tocBtn');
var toc = document.getElementById('toc');
var searchInput = document.getElementById('searchInput');
var progressBar = document.getElementById('progressBar');
var noResults = document.getElementById('noResults');
var docPage = document.querySelector('.docs-page');
var scrollTopBtn = document.getElementById('scrollTop');

var manualOpen = false;
var hideTimer = null;

function openToc(manual) {
  toc.classList.add('open');
  if (manual) { manualOpen = true; if (hideTimer) { clearTimeout(hideTimer); hideTimer = null; } }
}
function closeToc() {
  toc.classList.remove('open');
  manualOpen = false;
}
btn.addEventListener('click', function() { openToc(true); });
toc.querySelectorAll('a').forEach(function(a) { a.addEventListener('click', closeToc); });
document.addEventListener('click', function(e) {
  if (toc.classList.contains('open') && !toc.contains(e.target) && e.target !== btn && !btn.contains(e.target)) {
    closeToc();
  }
});

/* ----- progress bar + TOC auto-show ----- */
function onScroll() {
  var scrollTop = window.scrollY;
  var docHeight = document.documentElement.scrollHeight - window.innerHeight;
  var pct = docHeight > 0 ? Math.min(scrollTop / docHeight * 100, 100) : 0;
  progressBar.style.width = pct + '%';

  /* auto-show TOC on scroll, hide after idle */
  if (!manualOpen) {
    toc.classList.add('open');
    if (hideTimer) clearTimeout(hideTimer);
    hideTimer = setTimeout(function() {
      if (!manualOpen) closeToc();
      hideTimer = null;
    }, 2000);
  }

  /* active section */
  updateActiveSection();
}
window.addEventListener('scroll', onScroll);

/* ----- active TOC section ----- */
var headings = [];
docPage.querySelectorAll('h2[id], h3[id]').forEach(function(h) {
  headings.push({ el: h, id: h.id });
});

function updateActiveSection() {
  var scrollY = window.scrollY + 120;
  var current = null;
  for (var i = headings.length - 1; i >= 0; i--) {
    if (headings[i].el.offsetTop <= scrollY) { current = headings[i].id; break; }
  }
  if (!current && headings.length > 0) current = headings[0].id;

  toc.querySelectorAll('a').forEach(function(a) {
    a.classList.toggle('active', a.getAttribute('data-target') === current);
  });
}

/* ----- initial states ----- */
updateActiveSection();

/* ----- line numbers ----- */
var pres = document.querySelectorAll('.docs-page pre');
pres.forEach(function(pre) {
  var html = pre.innerHTML;
  var lines = html.split('\n');
  var wrapped = lines.map(function(line) {
    return '<span class="line">' + (line || ' ') + '</span>';
  }).join('\n');
  pre.innerHTML = wrapped;
  pre.classList.add('line-numbers');
});

/* ----- copy buttons ----- */
pres = document.querySelectorAll('.docs-page pre');
pres.forEach(function(pre) {
  var btn = document.createElement('button');
  btn.className = 'copy-btn';
  btn.textContent = 'copy';
  btn.setAttribute('aria-label', 'Copy code block');
  pre.appendChild(btn);
  btn.addEventListener('click', function() {
    var text = pre.textContent.trim();
    if (!text) return;
    if (navigator.clipboard) {
      navigator.clipboard.writeText(text).then(function() {
        btn.textContent = 'copied!';
        btn.classList.add('copied');
        setTimeout(function() { btn.textContent = 'copy'; btn.classList.remove('copied'); }, 2000);
      });
    } else {
      var ta = document.createElement('textarea');
      ta.value = text;
      ta.style.position = 'fixed';
      ta.style.opacity = '0';
      document.body.appendChild(ta);
      ta.select();
      document.execCommand('copy');
      document.body.removeChild(ta);
      btn.textContent = 'copied!';
      btn.classList.add('copied');
      setTimeout(function() { btn.textContent = 'copy'; btn.classList.remove('copied'); }, 2000);
    }
  });
});

/* ----- scroll to top ----- */
window.addEventListener('scroll', function() {
  scrollTopBtn.classList.toggle('visible', window.scrollY > 300);
});
scrollTopBtn.addEventListener('click', function() {
  window.scrollTo({ top: 0, behavior: 'smooth' });
});

/* ----- search ----- */
function getSearchableText(el) {
  var txt = '';
  el.childNodes.forEach(function(node) {
    if (node.nodeType === 3) txt += node.textContent;
    else if (node.nodeType === 1 && node.tagName !== 'SCRIPT' && node.tagName !== 'STYLE') {
      txt += getSearchableText(node);
    }
  });
  return txt;
}

var searchable = [];
docPage.querySelectorAll('h2, h3, p, pre, table, li, .syntax-diagram').forEach(function(el) {
  var text = getSearchableText(el).toLowerCase().replace(/\s+/g, ' ').trim();
  if (text) searchable.push({ el: el, text: text });
});

function doSearch(query) {
  var q = query.toLowerCase().replace(/\s+/g, ' ').trim();
  var all = docPage.querySelectorAll('h2, h3, p, pre, table, li, .syntax-diagram, ul');
  all.forEach(function(el) { el.classList.remove('search-match'); });
  if (!q) {
    docPage.classList.remove('search-active');
    noResults.style.display = 'none';
    return;
  }
  docPage.classList.add('search-active');
  var terms = q.split(' ').filter(Boolean);
  var matchCount = 0;
  searchable.forEach(function(item) {
    var match = true;
    for (var i = 0; i < terms.length; i++) {
      if (item.text.indexOf(terms[i]) === -1) { match = false; break; }
    }
    if (match) {
      item.el.classList.add('search-match');
      matchCount++;
    }
  });
  noResults.style.display = matchCount ? 'none' : 'block';
}

searchInput.addEventListener('input', function() { doSearch(this.value); });

document.addEventListener('keydown', function(e) {
  if (e.key === '/' && document.activeElement !== searchInput && document.activeElement.tagName !== 'INPUT' && document.activeElement.tagName !== 'TEXTAREA') {
    e.preventDefault();
    searchInput.focus();
  }
  if (e.key === 'Escape' && document.activeElement === searchInput) {
    searchInput.value = '';
    doSearch('');
    searchInput.blur();
  }
  if (e.key === 't' && !e.ctrlKey && !e.metaKey && document.activeElement === document.body) {
    if (toc.classList.contains('open')) closeToc(); else openToc(true);
  }
});
