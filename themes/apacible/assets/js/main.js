// Back to top button functionality with Intersection Observer
(function() {
  const backToTopButton = document.getElementById('back-to-top');
  const inlineBackToTop = document.querySelector('.back-to-top-link');

  // Scroll to top function - always scrolls to absolute top
  function scrollToTop(e) {
    e.preventDefault();
    e.stopPropagation();

    // Smooth scroll to top
    window.scrollTo({
      top: 0,
      left: 0,
      behavior: 'smooth'
    });
  }

  // Add event listeners
  if (backToTopButton) {
    backToTopButton.addEventListener('click', scrollToTop);

    // Use Intersection Observer instead of scroll events for better performance
    // Create a sentinel element at the top of the page
    const sentinel = document.createElement('div');
    sentinel.style.position = 'absolute';
    sentinel.style.top = '300px';
    sentinel.style.height = '1px';
    sentinel.style.width = '1px';
    sentinel.style.pointerEvents = 'none';
    document.body.prepend(sentinel);

    const observer = new IntersectionObserver(
      function(entries) {
        entries.forEach(function(entry) {
          // Show button when sentinel is out of view (scrolled past 300px)
          backToTopButton.classList.toggle('visible', !entry.isIntersecting);
        });
      },
      { threshold: 0 }
    );

    observer.observe(sentinel);
  }

  // Add event listener for inline link
  if (inlineBackToTop) {
    inlineBackToTop.addEventListener('click', scrollToTop);
  }
})();

// Theme toggle functionality
(function() {
  const themeToggles = document.querySelectorAll('[id="theme-toggle"]');
  if (!themeToggles.length) return;
  
  const htmlElement = document.documentElement;
  const THEME_KEY = 'theme';
  const DEFAULT_THEME = 'light';
  
  // Apply theme to all elements
  function applyTheme(theme) {
    htmlElement.setAttribute('data-theme', theme);
    document.body.setAttribute('data-theme', theme);
    // Update all theme toggle buttons if they exist
    themeToggles.forEach(function(toggle) {
      if (toggle) {
        toggle.setAttribute('data-theme', theme);
      }
    });
  }
  
  // Initialize theme from localStorage with error handling
  let currentTheme;
  try {
    currentTheme = localStorage.getItem(THEME_KEY) || DEFAULT_THEME;
  } catch (e) {
    console.error('Failed to read theme preference:', e);
    currentTheme = DEFAULT_THEME;
  }
  applyTheme(currentTheme);
  
  // Add click event to all theme toggle buttons
  themeToggles.forEach(function(themeToggle) {
    themeToggle.addEventListener('click', function() {
      currentTheme = currentTheme === 'light' ? 'dark' : 'light';
      
      // Save to localStorage with error handling
      try {
        localStorage.setItem(THEME_KEY, currentTheme);
      } catch (e) {
        console.error('Failed to save theme preference:', e);
      }
      
      applyTheme(currentTheme);
    });
  });
  
  // Re-sync theme when page is restored from bfcache (back/forward navigation)
  window.addEventListener('pageshow', function(event) {
    if (event.persisted) {
      // Page was restored from bfcache, re-read theme from localStorage
      try {
        currentTheme = localStorage.getItem(THEME_KEY) || DEFAULT_THEME;
        applyTheme(currentTheme);
      } catch (e) {
        console.error('Failed to sync theme on pageshow:', e);
      }
    }
  });
})();

// Toggle original/dithered images
function toggleOriginal(element) {
  const figure = element.closest('figure');
  if (!figure) return;
  
  const container = figure.querySelector('.image-container');
  if (!container) return;
  
  const dithered = container.querySelector('.dithered-img');
  const original = container.querySelector('.original-img');
  const showText = element.querySelector('.show-text');
  const hideText = element.querySelector('.hide-text');
  
  if (!dithered || !original) return;
  
  // Toggle active class
  const isActive = dithered.classList.contains('active');
  
  if (isActive) {
    // Show original
    dithered.classList.remove('active');
    dithered.style.display = 'none';
    original.style.display = 'block';
    if (showText) showText.style.display = 'none';
    if (hideText) hideText.style.display = 'inline';
  } else {
    // Show dithered
    dithered.classList.add('active');
    dithered.style.display = 'block';
    original.style.display = 'none';
    if (showText) showText.style.display = 'inline';
    if (hideText) hideText.style.display = 'none';
  }
}

// Cursor gradient effect - only on headings and links
(function() {
  // Only enable on desktop
  if (window.innerWidth <= 680) return;
  
  // Create cursor element
  const cursor = document.createElement('div');
  cursor.id = 'cursor';
  cursor.className = 'gradient';
  cursor.style.opacity = '0';
  document.body.appendChild(cursor);
  
  let mouseX = 0;
  let mouseY = 0;
  let cursorX = 0;
  let cursorY = 0;
  let isOverTarget = false;
  
  // Check if element is a heading or link
  function isTargetElement(element) {
    if (!element) return false;
    
    // Check if it's a heading
    if (/^H[1-6]$/.test(element.tagName)) return true;
    
    // Check if it's a link or inside a link
    if (element.tagName === 'A') return true;
    if (element.closest('a')) return true;
    
    return false;
  }
  
  // Track mouse position and hover state (throttled)
  let rafId = null;
  document.addEventListener('mousemove', function(e) {
    mouseX = e.clientX;
    mouseY = e.clientY;
    
    if (!rafId) {
      rafId = requestAnimationFrame(function() {
        const target = document.elementFromPoint(mouseX, mouseY);
        const shouldShow = isTargetElement(target);
        
        if (shouldShow !== isOverTarget) {
          isOverTarget = shouldShow;
          cursor.style.opacity = shouldShow ? '1' : '0';
        }
        rafId = null;
      });
    }
  }, { passive: true });
  
  // Smooth cursor movement
  function animateCursor() {
    const speed = 0.15;
    cursorX += (mouseX - cursorX) * speed;
    cursorY += (mouseY - cursorY) * speed;
    
    cursor.style.left = (cursorX - 125) + 'px';
    cursor.style.top = (cursorY - 125) + 'px';
    
    requestAnimationFrame(animateCursor);
  }
  
  animateCursor();
})();

// Page transitions
(function() {
  // Add fade-in on page load
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initFadeIn);
  } else {
    initFadeIn();
  }
  
  function initFadeIn() {
    document.body.style.opacity = '0';
    requestAnimationFrame(function() {
      document.body.style.transition = 'opacity 0.3s ease';
      document.body.style.opacity = '1';
    });
  }
  
  // Smooth transitions for internal links (if View Transitions API is supported)
  if (document.startViewTransition) {
    document.addEventListener('click', function(e) {
      const link = e.target.closest('a[href]');
      if (!link || link.target || !link.href || link.href.startsWith('#')) return;
      
      // Only for same-origin links
      try {
        const url = new URL(link.href);
        if (url.origin === window.location.origin) {
          e.preventDefault();
          document.startViewTransition(function() {
            window.location.href = link.href;
          });
        }
      } catch (err) {
        // Invalid URL, ignore
      }
    });
  }
})();

// Remove inline styles from code blocks (Zola adds background colors)
(function() {
  const STYLES_TO_REMOVE = ['background-color', 'background', 'color'];
  
  function removeInlineStyles(elements) {
    elements.forEach(function(el) {
      STYLES_TO_REMOVE.forEach(function(prop) {
        if (el.style[prop] || el.style.getPropertyValue(prop)) {
          el.style.removeProperty(prop);
        }
      });
    });
  }
  
  function processCodeBlocks() {
    removeInlineStyles(document.querySelectorAll('pre'));
    removeInlineStyles(document.querySelectorAll('code'));
  }
  
  // Run on page load
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', processCodeBlocks);
  } else {
    processCodeBlocks();
  }
  
  // Also run after dynamic content loads
  const observer = new MutationObserver(function(mutations) {
    const newCodeElements = [];
    mutations.forEach(function(mutation) {
      mutation.addedNodes.forEach(function(node) {
        if (node.nodeType === 1) {
          if (node.tagName === 'PRE' || node.tagName === 'CODE') {
            newCodeElements.push(node);
          }
          const nested = node.querySelectorAll && node.querySelectorAll('pre, code');
          if (nested) {
            newCodeElements.push(...Array.from(nested));
          }
        }
      });
    });
    
    if (newCodeElements.length > 0) {
      removeInlineStyles(newCodeElements);
    }
  });
  
  observer.observe(document.body, {
    childList: true,
    subtree: true
  });
})();

// Copy to clipboard functionality for code blocks
(function() {
  function addCopyButtons() {
    const codeBlocks = document.querySelectorAll('pre code');
    
    codeBlocks.forEach(function(codeBlock) {
      const pre = codeBlock.parentElement;
      
      // Skip if button already exists
      if (pre.querySelector('.copy-button')) {
        return;
      }
      
      // Create copy button
      const copyButton = document.createElement('button');
      copyButton.className = 'copy-button';
      copyButton.textContent = 'Copy';
      copyButton.setAttribute('aria-label', 'Copy code to clipboard');
      
      // Add click handler
      copyButton.addEventListener('click', function() {
        const text = codeBlock.textContent || codeBlock.innerText;
        const RESET_DELAY = 2000;
        
        function showFeedback(message, isError) {
          copyButton.textContent = message;
          if (!isError) {
            copyButton.classList.add('copied');
          }
          
          setTimeout(function() {
            copyButton.textContent = 'Copy';
            copyButton.classList.remove('copied');
          }, RESET_DELAY);
        }
        
        // Use modern Clipboard API if available
        if (navigator.clipboard && navigator.clipboard.writeText) {
          navigator.clipboard.writeText(text)
            .then(function() {
              showFeedback('Copied!', false);
            })
            .catch(function(err) {
              console.error('Failed to copy:', err);
              showFeedback('Error', true);
            });
        } else {
          // Fallback for older browsers
          const textArea = document.createElement('textarea');
          textArea.value = text;
          textArea.style.position = 'fixed';
          textArea.style.opacity = '0';
          textArea.style.pointerEvents = 'none';
          document.body.appendChild(textArea);
          textArea.select();
          
          try {
            const success = document.execCommand('copy');
            if (success) {
              showFeedback('Copied!', false);
            } else {
              showFeedback('Error', true);
            }
          } catch (err) {
            console.error('Failed to copy:', err);
            showFeedback('Error', true);
          }
          
          document.body.removeChild(textArea);
        }
      });
      
      // Insert button into pre element
      pre.appendChild(copyButton);
    });
  }
  
  // Run on page load
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', addCopyButtons);
  } else {
    addCopyButtons();
  }
  
  // Also run after dynamic content loads (if using a framework)
  // This is useful if code blocks are added dynamically
  const observer = new MutationObserver(function(mutations) {
    let shouldAddButtons = false;
    mutations.forEach(function(mutation) {
      if (mutation.addedNodes.length > 0) {
        mutation.addedNodes.forEach(function(node) {
          if (node.nodeType === 1 && (node.tagName === 'PRE' || node.querySelector('pre'))) {
            shouldAddButtons = true;
          }
        });
      }
    });
    
    if (shouldAddButtons) {
      addCopyButtons();
    }
  });
  
  observer.observe(document.body, {
    childList: true,
    subtree: true
  });
})();

// Table of Contents - Substack-style with line indicators
(function() {
  const tocElement = document.querySelector('.toc');
  if (!tocElement) return;
  
  const tocLinks = document.querySelectorAll('.toc-content a');
  if (tocLinks.length === 0) return;
  
  const headings = Array.from(document.querySelectorAll('h1[id], h2[id], h3[id], h4[id], h5[id], h6[id]'));
  
  // Create lines container and generate lines based on TOC structure
  function initializeTOCLines() {
    // Check if we're on large screen
    if (window.innerWidth < 940) return;
    
    // Remove existing lines if any
    const existingLines = tocElement.querySelector('.toc-lines');
    if (existingLines) existingLines.remove();
    
    // Create lines container
    const linesContainer = document.createElement('div');
    linesContainer.className = 'toc-lines';
    
    // Generate a line for each TOC link
    tocLinks.forEach(function(link) {
      const href = link.getAttribute('href');
      if (!href || !href.startsWith('#')) return;
      
      const targetId = href.substring(1);
      const targetHeading = document.getElementById(targetId);
      if (!targetHeading) return;
      
      // Determine heading level
      const tagName = targetHeading.tagName.toLowerCase();
      const level = tagName.charAt(1); // Extract number from h1, h2, etc.
      
      // Create line element
      const line = document.createElement('div');
      line.className = 'toc-line';
      line.setAttribute('data-level', Math.min(level, 3)); // Cap at level 3
      line.setAttribute('data-target', targetId);
      
      linesContainer.appendChild(line);
    });
    
    // Add click handler to lines container to toggle TOC
    linesContainer.addEventListener('click', function(e) {
      e.stopPropagation();
      toggleTOC();
    });
    
    // Insert lines before toc-content
    const tocContent = tocElement.querySelector('.toc-content');
    if (tocContent) {
      tocElement.insertBefore(linesContainer, tocContent);
    }
  }
  
  // Toggle TOC expanded/collapsed state
  function toggleTOC() {
    tocElement.classList.toggle('toc-expanded');
  }
  
  // Close TOC when clicking outside
  function closeTOC() {
    tocElement.classList.remove('toc-expanded');
  }
  
  // Smooth scroll to heading with animation
  function scrollToHeading(targetId) {
    const target = document.getElementById(targetId);
    if (!target) return;
    
    const offset = 80;
    const targetPosition = target.getBoundingClientRect().top + window.pageYOffset - offset;
    
    window.scrollTo({
      top: targetPosition,
      behavior: 'smooth'
    });
  }
  
  // Update active line and link highlighting
  function updateActiveTOC() {
    let current = '';
    
    // Find the heading currently in view
    headings.forEach(function(heading) {
      const rect = heading.getBoundingClientRect();
      if (rect.top <= 100) {
        current = heading.id;
      }
    });
    
    // Update active state on links
    tocLinks.forEach(function(link) {
      const href = link.getAttribute('href');
      const id = href ? href.substring(1) : '';
      
      if (id === current) {
        link.classList.add('active');
      } else {
        link.classList.remove('active');
      }
    });
    
    // Update active state on lines
    const lines = document.querySelectorAll('.toc-line');
    lines.forEach(function(line) {
      const targetId = line.getAttribute('data-target');
      
      if (targetId === current) {
        line.classList.add('active');
      } else {
        line.classList.remove('active');
      }
    });
  }
  
  // Initialize lines
  initializeTOCLines();
  
  // Update on scroll (throttled with requestAnimationFrame)
  let ticking = false;
  window.addEventListener('scroll', function() {
    if (!ticking) {
      requestAnimationFrame(function() {
        updateActiveTOC();
        ticking = false;
      });
      ticking = true;
    }
  }, { passive: true });
  
  // Update on page load
  updateActiveTOC();
  
  // Reinitialize on resize (in case screen size changes)
  let resizeTimeout;
  window.addEventListener('resize', function() {
    clearTimeout(resizeTimeout);
    resizeTimeout = setTimeout(function() {
      initializeTOCLines();
      updateActiveTOC();
    }, 250);
  });
  
  // Smooth scroll for TOC links in the expanded card
  tocLinks.forEach(function(link) {
    link.addEventListener('click', function(e) {
      const href = this.getAttribute('href');
      if (href && href.startsWith('#')) {
        const targetId = href.substring(1);
        e.preventDefault();
        scrollToHeading(targetId);
        // Close TOC after navigation
        closeTOC();
      }
    });
  });
  
  // Close TOC when clicking outside
  document.addEventListener('click', function(e) {
    if (!tocElement.contains(e.target) && tocElement.classList.contains('toc-expanded')) {
      closeTOC();
    }
  });
  
  // Prevent clicks inside TOC card from closing it
  const tocContent = tocElement.querySelector('.toc-content');
  if (tocContent) {
    tocContent.addEventListener('click', function(e) {
      e.stopPropagation();
    });
  }
})();

// Prefetch recent post links on hover for faster navigation
(function() {
  const recentPostsSection = document.querySelector('.recent-posts');
  if (!recentPostsSection) return;
  
  const links = recentPostsSection.querySelectorAll('a[href]');
  const prefetchedUrls = new Set();
  const HOVER_DELAY = 100;
  
  function prefetchUrl(url) {
    if (prefetchedUrls.has(url)) return;
    
    const linkElement = document.createElement('link');
    linkElement.rel = 'prefetch';
    linkElement.href = url;
    linkElement.as = 'document';
    document.head.appendChild(linkElement);
    prefetchedUrls.add(url);
  }
  
  links.forEach(function(link) {
    const url = link.getAttribute('href');
    if (!url || url.startsWith('#')) return;
    
    // Prefetch on hover with a small delay to avoid unnecessary requests
    let hoverTimeout;
    link.addEventListener('mouseenter', function() {
      hoverTimeout = setTimeout(function() {
        prefetchUrl(url);
      }, HOVER_DELAY);
    }, { passive: true });
    
    link.addEventListener('mouseleave', function() {
      if (hoverTimeout) {
        clearTimeout(hoverTimeout);
      }
    }, { passive: true });
  });
})();

