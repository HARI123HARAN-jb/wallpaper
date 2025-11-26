document.addEventListener('DOMContentLoaded', () => {
    const dropZone = document.getElementById('drop-zone');
    const fileInput = document.getElementById('file-input');
    const generateBtn = document.getElementById('generate-btn');
    const progressContainer = document.getElementById('progress-container');
    const resultsArea = document.getElementById('results-area');
    const gallery = document.getElementById('gallery');
    const browseBtn = document.querySelector('.browse-btn');

    let selectedFile = null;

    // Drag & Drop Handlers
    dropZone.addEventListener('dragover', (e) => {
        e.preventDefault();
        dropZone.classList.add('dragover');
    });

    dropZone.addEventListener('dragleave', () => {
        dropZone.classList.remove('dragover');
    });

    dropZone.addEventListener('drop', (e) => {
        e.preventDefault();
        dropZone.classList.remove('dragover');
        if (e.dataTransfer.files.length) {
            handleFileSelect(e.dataTransfer.files[0]);
        }
    });

    // Click to browse
    dropZone.addEventListener('click', () => fileInput.click());
    
    fileInput.addEventListener('change', (e) => {
        if (e.target.files.length) {
            handleFileSelect(e.target.files[0]);
        }
    });

    function handleFileSelect(file) {
        if (!file.type.startsWith('image/')) {
            alert('Please upload an image file.');
            return;
        }
        selectedFile = file;
        dropZone.querySelector('h3').textContent = file.name;
        dropZone.querySelector('p').textContent = 'Ready to generate';
    }

    // Generate Button Handler
    generateBtn.addEventListener('click', async () => {
        if (!selectedFile) {
            alert('Please select an image first.');
            return;
        }

        const checkboxes = document.querySelectorAll('input[name="res"]:checked');
        if (checkboxes.length === 0) {
            alert('Please select at least one resolution.');
            return;
        }

        const resolutions = Array.from(checkboxes).map(cb => cb.value).join(',');

        // UI Updates
        generateBtn.disabled = true;
        generateBtn.textContent = 'Processing...';
        progressContainer.classList.remove('hidden');
        resultsArea.classList.add('hidden');
        gallery.innerHTML = '';

        // Prepare FormData
        const formData = new FormData();
        formData.append('image', selectedFile);
        formData.append('resolutions', resolutions);

        try {
            const response = await fetch('/upload', {
                method: 'POST',
                body: formData
            });

            if (!response.ok) {
                throw new Error('Upload failed');
            }

            const data = await response.json();
            displayResults(data.images);

        } catch (error) {
            console.error('Error:', error);
            alert('An error occurred while processing the image.');
        } finally {
            generateBtn.disabled = false;
            generateBtn.textContent = 'Generate Wallpapers';
            progressContainer.classList.add('hidden');
        }
    });

    function displayResults(images) {
        resultsArea.classList.remove('hidden');
        
        images.forEach(img => {
            const card = document.createElement('div');
            card.className = 'result-card';
            
            card.innerHTML = `
                <img src="${img.url}" class="result-img" alt="${img.res}">
                <div class="result-info">
                    <span class="result-res">${img.res}</span>
                    <a href="${img.url}" download class="download-link">Download</a>
                </div>
            `;
            
            gallery.appendChild(card);
        });
        
        // Scroll to results
        resultsArea.scrollIntoView({ behavior: 'smooth' });
    }
});
