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
    // Check if Alpine.js is handling filters (check for x-data attribute on collection-wrapper)
    const alpineWrapper = document.querySelector('[x-data]');
    const isAlpineControlled = alpineWrapper && alpineWrapper.querySelector('.collection-filter');

    // Only initialize vanilla JS filters if Alpine.js is NOT controlling them
    if (!isAlpineControlled) {
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
    }

    // Initialize poster detail view (always needed for tooltips)
    initPosterDetails();

    // Spoiler toggles are now handled by Alpine.js - no need to initialize here
  }

  function initYearFilter(filterContainer) {
    const buttons = filterContainer.querySelectorAll('.filter-btn');
    // Find the collection container - it's after the .collection-filters wrapper
    let collection = filterContainer.closest('.collection-filters');
    if (collection) {
      collection = collection.nextElementSibling;
      while (collection && !collection.classList.contains('collection')) {
        collection = collection.nextElementSibling;
      }
    } else {
      // Fallback: try nextElementSibling approach
      collection = filterContainer.nextElementSibling;
      while (collection && !collection.classList.contains('collection')) {
        collection = collection.nextElementSibling;
      }
    }
    
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
    // Find the collection container - it's after the .collection-filters wrapper
    let collection = filterContainer.closest('.collection-filters');
    if (collection) {
      collection = collection.nextElementSibling;
      while (collection && !collection.classList.contains('collection')) {
        collection = collection.nextElementSibling;
      }
    } else {
      // Fallback: try nextElementSibling approach
      collection = filterContainer.nextElementSibling;
      while (collection && !collection.classList.contains('collection')) {
        collection = collection.nextElementSibling;
      }
    }
    
    if (!collection || !collection.dataset.filterable) return;

    buttons.forEach(btn => {
      btn.addEventListener('click', function() {
        const selectedYear = this.dataset.date;
        
        // Update active state
        buttons.forEach(b => b.classList.remove('active'));
        this.classList.add('active');
        
        // Filter items by year consumed
        // data-date can be YYYY (for year filter) or YYYY-MM-DD/YYYY-MM (for date_read)
        const items = collection.querySelectorAll('.collection-item[data-date]');
        items.forEach(item => {
          const itemDate = item.dataset.date;
          let itemYear = '';
          
          // Extract year from data-date
          if (itemDate) {
            // If it's just YYYY (4 digits), use it directly
            if (/^\d{4}$/.test(itemDate)) {
              itemYear = itemDate;
            } else {
              // Otherwise, extract year from YYYY-MM-DD or YYYY-MM format
              itemYear = itemDate.split('-')[0];
            }
          }
          
          if (selectedYear === 'all' || itemYear === selectedYear) {
            item.style.display = '';
            item.removeAttribute('aria-hidden');
          } else {
            item.style.display = 'none';
            item.setAttribute('aria-hidden', 'true');
          }
        });

        // Also respect category filter if active
        const filtersWrapper = filterContainer.closest('.collection-filters');
        if (filtersWrapper) {
          const categoryFilter = filtersWrapper.querySelector('.collection-filter[data-filter-type="category"]');
          if (categoryFilter) {
            const activeCategoryBtn = categoryFilter.querySelector('.filter-btn.active');
            if (activeCategoryBtn && activeCategoryBtn.dataset.category !== 'all') {
              const selectedCategory = activeCategoryBtn.dataset.category;
              items.forEach(item => {
                if (item.style.display !== 'none') {
                  if (item.dataset.category !== selectedCategory) {
                    item.style.display = 'none';
                    item.setAttribute('aria-hidden', 'true');
                  }
                }
              });
            }
          }
        }

        // Announce change for screen readers
        const visibleCount = collection.querySelectorAll('.collection-item[data-date]:not([aria-hidden])').length;
        announceYearChange(selectedYear, visibleCount);
      });
    });
  }

  function initCategoryFilter(filterContainer) {
    const buttons = filterContainer.querySelectorAll('.filter-btn');
    // Find the collection container - it's after the .collection-filters wrapper
    let collection = filterContainer.closest('.collection-filters');
    if (collection) {
      collection = collection.nextElementSibling;
      while (collection && !collection.classList.contains('collection')) {
        collection = collection.nextElementSibling;
      }
    } else {
      // Fallback: try nextElementSibling approach
      collection = filterContainer.nextElementSibling;
      while (collection && !collection.classList.contains('collection')) {
        collection = collection.nextElementSibling;
      }
    }
    
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

        // Also respect year filter if active
        const filtersWrapper = filterContainer.closest('.collection-filters');
        if (filtersWrapper) {
          const yearFilter = filtersWrapper.querySelector('.collection-filter[data-filter-type="date"]');
          if (yearFilter) {
            const activeYearBtn = yearFilter.querySelector('.filter-btn.active');
            if (activeYearBtn && activeYearBtn.dataset.date !== 'all') {
              const selectedYear = activeYearBtn.dataset.date;
              items.forEach(item => {
                if (item.style.display !== 'none') {
                  const itemDate = item.dataset.date;
                  let itemYear = '';
                  if (itemDate) {
                    if (/^\d{4}$/.test(itemDate)) {
                      itemYear = itemDate;
                    } else {
                      itemYear = itemDate.split('-')[0];
                    }
                  }
                  if (itemYear !== selectedYear) {
                    item.style.display = 'none';
                    item.setAttribute('aria-hidden', 'true');
                  }
                }
              });
            }
          }
        }

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

  function announceYearChange(year, visibleCount) {
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
      ? `Showing all items (${visibleCount})` 
      : `Filtered by year ${year} (${visibleCount} items)`;
    liveRegion.textContent = message;
  }

  function initPosterDetails() {
    // Find all poster grid collections
    const posterCollections = document.querySelectorAll('.poster-grid');
    if (!posterCollections.length) return;

    // Process each poster collection
    posterCollections.forEach(posterCollection => {
      const posterItems = posterCollection.querySelectorAll('.collection-item');
      posterItems.forEach(item => {
        const img = item.querySelector('img');
        if (!img) return; // Skip if no image
        
        // Get data from attributes
        const data = {
          title: item.dataset.title || '',
          year: item.dataset.year || '',
          artist: item.dataset.artist || '',
          category: item.dataset.category || ''
        };

        // Create tooltip text based on category
        let tooltipText = '';
        if (data.category === 'movies') {
          tooltipText = data.title;
          if (data.year) {
            tooltipText += ` (${data.year})`;
          }
        } else if (data.category === 'music') {
          if (data.artist && data.title) {
            tooltipText = `${data.artist} - ${data.title}`;
          } else if (data.title) {
            tooltipText = data.title;
          }
          if (data.year) {
            tooltipText += ` (${data.year})`;
          }
        } else if (data.category === 'books') {
          tooltipText = data.title;
          if (data.year) {
            tooltipText += ` (${data.year})`;
          }
        } else {
          tooltipText = data.title;
        }

        // Set tooltip on the image or link
        if (tooltipText) {
          const link = img.closest('a');
          if (link) {
            link.setAttribute('title', tooltipText);
          } else {
            img.setAttribute('title', tooltipText);
          }
        }
      });
    });
  }

  // Spoiler toggles are now handled by Alpine.js - removed
})();

