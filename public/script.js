(function(){
  'use strict';

  /* ---- Particle background ---- */
  const canvas = document.getElementById('bg');
  const ctx = canvas.getContext('2d');
  let particles = [];
  let mouse = { x: 0, y: 0 };
  let animId;

  function resize() {
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
  }

  window.addEventListener('resize', resize);
  document.addEventListener('mousemove', e => { mouse.x = e.clientX; mouse.y = e.clientY; });

  class Particle {
    constructor() { this.reset(); }
    reset() {
      this.x = Math.random() * canvas.width;
      this.y = Math.random() * canvas.height;
      this.size = 1 + Math.random() * 2;
      this.speedX = (Math.random() - .5) * .3;
      this.speedY = (Math.random() - .5) * .3;
      this.opacity = .2 + Math.random() * .3;
    }
    update() {
      this.x += this.speedX;
      this.y += this.speedY;
      if (this.x < 0 || this.x > canvas.width) this.speedX *= -1;
      if (this.y < 0 || this.y > canvas.height) this.speedY *= -1;
    }
    draw() {
      ctx.beginPath();
      ctx.arc(this.x, this.y, this.size, 0, Math.PI * 2);
      ctx.fillStyle = `rgba(46, 160, 67, ${this.opacity})`;
      ctx.fill();
    }
  }

  function initParticles() {
    const count = Math.min(80, Math.floor(canvas.width * canvas.height / 12000));
    particles = Array.from({ length: count }, () => new Particle());
  }

  function drawLines() {
    for (let i = 0; i < particles.length; i++) {
      for (let j = i + 1; j < particles.length; j++) {
        const dx = particles[i].x - particles[j].x;
        const dy = particles[i].y - particles[j].y;
        const dist = Math.sqrt(dx * dx + dy * dy);
        if (dist < 140) {
          const opacity = (1 - dist / 140) * .12;
          ctx.beginPath();
          ctx.moveTo(particles[i].x, particles[i].y);
          ctx.lineTo(particles[j].x, particles[j].y);
          ctx.strokeStyle = `rgba(46, 160, 67, ${opacity})`;
          ctx.lineWidth = .5;
          ctx.stroke();
        }
      }
    }
  }

  function animate() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    particles.forEach(p => { p.update(); p.draw(); });
    drawLines();
    animId = requestAnimationFrame(animate);
  }

  resize();
  initParticles();
  animate();

  /* ---- Typewriter code demo with syntax highlighting ---- */
  const code = `var.name "World"

fn.greet {
  print $"Hello, {name}!"
}

run.greet

var.fruits ("apple", "banana")
for fruit in fruits {
  print fruit
}`;

  const el = document.getElementById('demo-code');
  let idx = 0;
  let isDeleting = false;
  let pauseType = ''; // '' = running, 'done' = full pause, 'empty' = reset pause
  let lastDomUpdate = 0;

  const kw = new Set(['print','input','for','in','if','elif','else','True','False','Null']);
  const prefixes = new Set(['var','fn','run','add']);

  function esc(s) {
    return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
  }

  function highlightSource(src) {
    const out = [];
    let i = 0;
    while (i < src.length) {
      if (src[i] === '!' && (i === 0 || src[i-1] === '\n')) {
        let end = src.indexOf('\n', i);
        if (end === -1) end = src.length;
        out.push('<span class="cm">' + esc(src.slice(i, end)) + '</span>');
        i = end;
        continue;
      }
      if (src[i] === '$' && src[i+1] === '"') {
        let end = src.indexOf('"', i + 2);
        if (end === -1) end = src.length;
        let inner = src.slice(i + 2, end);
        inner = inner.replace(/\{([^}]+)\}/g, function(_, expr) {
          return '<span class="fn">{' + esc(expr) + '}</span>';
        });
        out.push('<span class="str">$&quot;' + inner + '&quot;</span>');
        i = end + 1;
        continue;
      }
      if (src[i] === '"') {
        let end = src.indexOf('"', i + 1);
        if (end === -1) end = src.length - 1;
        out.push('<span class="str">' + esc(src.slice(i, end + 1)) + '</span>');
        i = end + 1;
        continue;
      }
      let rest = src.slice(i);
      if (rest.startsWith('range:')) {
        out.push('<span class="kw">range</span>:');
        i += 6;
        continue;
      }
      let m = rest.match(/^(var|fn|run|add)\./);
      if (m) {
        out.push('<span class="kw">' + m[1] + '</span>.');
        i += m[0].length;
        continue;
      }
      m = rest.match(/^([a-zA-Z_]\w*)/);
      if (m && kw.has(m[1])) {
        out.push('<span class="kw">' + m[1] + '</span>');
        i += m[1].length;
        continue;
      }
      if (src[i] === '.') {
        let idMatch = src.slice(i + 1).match(/^[a-zA-Z_]\w*/);
        if (idMatch) {
          out.push('.<span class="fn">' + idMatch[0] + '</span>');
          i += 1 + idMatch[0].length;
          continue;
        }
      }
      m = rest.match(/^\d+(\.\d+)?/);
      if (m) {
        out.push('<span class="num">' + m[0] + '</span>');
        i += m[0].length;
        continue;
      }
      out.push(esc(src[i]));
      i++;
    }
    return out.join('');
  }

  function flushDisplay() {
    el.innerHTML = idx > 0 ? highlightSource(code.slice(0, idx)) : '';
    lastDomUpdate = performance.now();
  }

  function typewriter() {
    if (pauseType) {
      if (pauseType === 'done') {
        flushDisplay();
        pauseType = '';
        setTimeout(typewriter, 3000);
      } else {
        el.innerHTML = '';
        pauseType = '';
        setTimeout(typewriter, 800);
      }
      return;
    }

    if (!isDeleting) {
      idx++;
      if (idx >= code.length) {
        isDeleting = true;
        flushDisplay();
        pauseType = 'done';
      } else if (performance.now() - lastDomUpdate > 80) {
        flushDisplay();
      }
      const speed = idx < 20 ? 45 : 25 + Math.random() * 20;
      setTimeout(typewriter, speed);
    } else {
      idx--;
      if (idx <= 0) {
        isDeleting = false;
        el.innerHTML = '';
        pauseType = 'empty';
      } else if (performance.now() - lastDomUpdate > 70) {
        flushDisplay();
      }
      setTimeout(typewriter, 12 + Math.random() * 8);
    }
  }

  setTimeout(typewriter, 1400);

  /* ---- Dynamic version & file size from GitHub ---- */
  async function fetchReleaseInfo() {
    try {
      const resp = await fetch('https://api.github.com/repos/hannogeo/emerald/releases/latest');
      if (!resp.ok) return;
      const data = await resp.json();
      const tag = data.tag_name;
      const versionEl = document.getElementById('version-badge');
      if (versionEl) versionEl.textContent = tag;

      const asset = data.assets.find(a => a.name === 'emerald-installer.exe');
      if (asset) {
        const sizeEl = document.getElementById('file-size');
        if (sizeEl) {
          const mb = Math.round(asset.size / (1024 * 1024));
          sizeEl.textContent = '~' + mb + ' MB';
        }
        const dlBtn = document.getElementById('download-btn');
        if (dlBtn) dlBtn.href = asset.browser_download_url;
        const heroBtn = document.getElementById('hero-download-btn');
        if (heroBtn) heroBtn.href = '#download';
      }
    } catch (_) {}
  }

  fetchReleaseInfo();

  document.addEventListener('click', function(e) {
    var link = e.target.closest('a');
    if (!link) return;
    var href = link.getAttribute('href');
    if (!href || href.startsWith('http') || href.startsWith('#') || href.startsWith('//')) return;
    if (href === location.pathname.replace(/\/?$/, '/index.html') || href === location.pathname) return;
    e.preventDefault();
    window.location = href;
  });

})();
