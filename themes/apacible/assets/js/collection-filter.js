/**
 * Collection Filter - Minimal vanilla JS for filtering collections
 * Smolweb compliant: lightweight, accessible, progressive enhancement
 */

(function() {
  'use strict';

  // Initialize on DOM ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }

  function init() {
    const filters = document.querySelectorAll('.collection-filter');
    if (filters.length) {
    filters.forEach(filter => {
      const filterType = filter.dataset.filterType;
      if (filterType === 'year') {
        initYearFilter(filter);
      } else if (filterType === 'date') {
        initDateFilter(filter);
        } else if (filterType === 'category') {
          initCategoryFilter(filter);
      }
    });
    }
    
    // Initialize poster detail view
    initPosterDetails();
    
    // Initialize spoiler toggles
    initSpoilerToggles();
  }

  function initYearFilter(filterContainer) {
    const buttons = filterContainer.querySelectorAll('.filter-btn');
    const collection = filterContainer.nextElementSibling;
    
    if (!collection || !collection.dataset.filterable) return;

    buttons.forEach(btn => {
      btn.addEventListener('click', function() {
        const selectedYear = this.dataset.year;
        
        // Update active state
        buttons.forEach(b => b.classList.remove('active'));
        this.classList.add('active');
        
        // Filter items
        const items = collection.querySelectorAll('.collection-item[data-year]');
        items.forEach(item => {
          if (selectedYear === 'all' || item.dataset.year === selectedYear) {
            item.style.display = '';
            item.removeAttribute('aria-hidden');
          } else {
            item.style.display = 'none';
            item.setAttribute('aria-hidden', 'true');
          }
        });

        // Announce change for screen readers
        announceFilterChange(selectedYear, items.length);
      });
    });
  }

  function initDateFilter(filterContainer) {
    const buttons = filterContainer.querySelectorAll('.filter-btn');
    const collection = filterContainer.nextElementSibling;
    
    if (!collection || !collection.dataset.filterable) return;

    buttons.forEach(btn => {
      btn.addEventListener('click', function() {
        const selectedDate = this.dataset.date;
        
        // Update active state
        buttons.forEach(b => b.classList.remove('active'));
        this.classList.add('active');
        
        // Filter items by year (data-date format: YYYY-MM, button has YYYY)
        const items = collection.querySelectorAll('.collection-item[data-date]');
        items.forEach(item => {
          const itemYear = item.dataset.date.split('-')[0]; // Extract year from YYYY-MM
          if (selectedDate === 'all' || itemYear === selectedDate) {
            item.style.display = '';
            item.removeAttribute('aria-hidden');
          } else {
            item.style.display = 'none';
            item.setAttribute('aria-hidden', 'true');
          }
        });

        // Announce change for screen readers
        announceFilterChange(selectedDate, items.length);
      });
    });
  }

  function initCategoryFilter(filterContainer) {
    const buttons = filterContainer.querySelectorAll('.filter-btn');
    const collection = filterContainer.nextElementSibling;
    
    if (!collection || !collection.dataset.filterable) return;

    buttons.forEach(btn => {
      btn.addEventListener('click', function() {
        const selectedCategory = this.dataset.category;
        
        // Update active state
        buttons.forEach(b => b.classList.remove('active'));
        this.classList.add('active');
        
        // Filter items
        const items = collection.querySelectorAll('.collection-item[data-category]');
        items.forEach(item => {
          if (selectedCategory === 'all' || item.dataset.category === selectedCategory) {
            item.style.display = '';
            item.removeAttribute('aria-hidden');
          } else {
            item.style.display = 'none';
            item.setAttribute('aria-hidden', 'true');
          }
        });

        // Announce change for screen readers
        const visibleCount = collection.querySelectorAll('.collection-item[data-category]:not([aria-hidden])').length;
        announceCategoryChange(selectedCategory, visibleCount);
      });
    });
  }

  function announceFilterChange(year, totalItems) {
    // Create or update aria-live region for accessibility
    let liveRegion = document.getElementById('filter-announcement');
    if (!liveRegion) {
      liveRegion = document.createElement('div');
      liveRegion.id = 'filter-announcement';
      liveRegion.setAttribute('aria-live', 'polite');
      liveRegion.setAttribute('aria-atomic', 'true');
      liveRegion.className = 'sr-only';
      document.body.appendChild(liveRegion);
    }
    
    const message = year === 'all' 
      ? `Showing all items` 
      : `Filtered by ${year}`;
    liveRegion.textContent = message;
  }

  function announceCategoryChange(category, visibleCount) {
    // Create or update aria-live region for accessibility
    let liveRegion = document.getElementById('filter-announcement');
    if (!liveRegion) {
      liveRegion = document.createElement('div');
      liveRegion.id = 'filter-announcement';
      liveRegion.setAttribute('aria-live', 'polite');
      liveRegion.setAttribute('aria-atomic', 'true');
      liveRegion.className = 'sr-only';
      document.body.appendChild(liveRegion);
    }
    
    const categoryName = category.charAt(0).toUpperCase() + category.slice(1);
    const message = category === 'all' 
      ? `Showing all items (${visibleCount})` 
      : `Showing ${categoryName} (${visibleCount} items)`;
    liveRegion.textContent = message;
  }

  function initPosterDetails() {
    // Only activate for poster grid collections
    const posterCollection = document.querySelector('.collection-poster.poster-grid');
    if (!posterCollection) return;

    // Create overlay element
    const overlay = document.createElement('div');
    overlay.className = 'poster-detail-overlay';
    overlay.innerHTML = `
      <div class="poster-detail-content">
        <div class="poster-detail-header">
          <h3></h3>
          <button class="poster-detail-close" aria-label="Close">&times;</button>
        </div>
        <div class="poster-detail-meta"></div>
        <p class="poster-detail-description"></p>
        <div class="poster-detail-footer"></div>
      </div>
    `;
    document.body.appendChild(overlay);

    const content = overlay.querySelector('.poster-detail-content');
    const closeBtn = overlay.querySelector('.poster-detail-close');
    const title = overlay.querySelector('.poster-detail-header h3');
    const meta = overlay.querySelector('.poster-detail-meta');
    const description = overlay.querySelector('.poster-detail-description');
    const footer = overlay.querySelector('.poster-detail-footer');

    // Close overlay when clicking outside or close button
    overlay.addEventListener('click', (e) => {
      if (e.target === overlay) closeOverlay();
    });
    closeBtn.addEventListener('click', closeOverlay);

    // ESC key to close
    document.addEventListener('keydown', (e) => {
      if (e.key === 'Escape' && overlay.classList.contains('active')) {
        closeOverlay();
      }
    });

    function closeOverlay() {
      overlay.classList.remove('active');
      document.querySelectorAll('.collection-item.selected').forEach(item => {
        item.classList.remove('selected');
      });
    }

    // Add click handlers to poster items
    const posterItems = posterCollection.querySelectorAll('.collection-item');
    posterItems.forEach(item => {
      item.addEventListener('click', (e) => {
        e.preventDefault();
        
        // Highlight selected poster
        document.querySelectorAll('.collection-item.selected').forEach(i => {
          i.classList.remove('selected');
        });
        item.classList.add('selected');

        // Get data from attributes
        const data = {
          title: item.dataset.title || '',
          year: item.dataset.year || '',
          director: item.dataset.director || '',
          author: item.dataset.author || '',
          artist: item.dataset.artist || '',
          label: item.dataset.label || '',
          rating: item.dataset.rating || '',
          content: item.dataset.content || '',
          footer: item.dataset.footer || '',
          link: item.dataset.link || '',
          category: item.dataset.category || ''
        };

        // Populate overlay
        if (data.link && !data.link.startsWith('http')) {
          // Internal link - make title clickable
          title.innerHTML = `<a href="${data.link}">${data.title}</a>`;
        } else {
          title.textContent = data.title;
        }
        
        // Build meta information
        let metaHTML = '';
        if (data.director) {
          metaHTML += `<div class="poster-detail-meta-row"><span class="poster-detail-meta-label">Director:</span><span>${data.director}</span></div>`;
        }
        if (data.author) {
          metaHTML += `<div class="poster-detail-meta-row"><span class="poster-detail-meta-label">Author:</span><span>${data.author}</span></div>`;
        }
        if (data.artist) {
          metaHTML += `<div class="poster-detail-meta-row"><span class="poster-detail-meta-label">Artist:</span><span>${data.artist}</span></div>`;
        }
        if (data.label) {
          metaHTML += `<div class="poster-detail-meta-row"><span class="poster-detail-meta-label">Label:</span><span>${data.label}</span></div>`;
        }
        if (data.year) {
          metaHTML += `<div class="poster-detail-meta-row"><span class="poster-detail-meta-label">Year:</span><span>${data.year}</span></div>`;
        }
        if (data.rating) {
          const stars = '★'.repeat(parseInt(data.rating)) + '☆'.repeat(5 - parseInt(data.rating));
          metaHTML += `<div class="poster-detail-meta-row"><span class="poster-detail-meta-label">Rating:</span><span class="poster-detail-rating">${stars}</span></div>`;
        }
        meta.innerHTML = metaHTML;

        description.textContent = data.content;
        footer.textContent = data.footer;

        // Show overlay
        overlay.classList.add('active');
      });
    });
  }

  function initSpoilerToggles() {
    const spoilerToggles = document.querySelectorAll('.spoiler-toggle');
    if (!spoilerToggles.length) return;

    spoilerToggles.forEach(toggle => {
      const content = toggle.nextElementSibling;
      const showText = toggle.dataset.showText;
      const hideText = toggle.dataset.hideText;

      toggle.addEventListener('click', () => {
        const isHidden = content.hasAttribute('hidden');
        
        if (isHidden) {
          content.removeAttribute('hidden');
          toggle.classList.add('active');
          toggle.textContent = hideText;
          toggle.setAttribute('aria-expanded', 'true');
        } else {
          content.setAttribute('hidden', '');
          toggle.classList.remove('active');
          toggle.textContent = showText;
          toggle.setAttribute('aria-expanded', 'false');
        }
      });

      // Set initial aria-expanded
      toggle.setAttribute('aria-expanded', 'false');
    });
  }
})();

