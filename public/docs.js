const btn = document.getElementById('tocBtn');
const toc = document.getElementById('toc');
const overlay = document.getElementById('tocOverlay');

function openToc() {
  toc.classList.add('open');
  overlay.classList.add('open');
}

function closeToc() {
  toc.classList.remove('open');
  overlay.classList.remove('open');
}

btn.addEventListener('click', openToc);
overlay.addEventListener('click', closeToc);

toc.querySelectorAll('a').forEach(function(a) {
  a.addEventListener('click', closeToc);
});
