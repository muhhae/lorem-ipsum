document.addEventListener('DOMContentLoaded', () => {
    const maxFileSize = 1024 * 1024 * 5
    const maxFileCount = 8
    // const allowedFileTypes = ['image/png', 'image/jpeg', 'image/gif', 'image/bmp', 'image/webp']
    var Images = []
    var selectedImages = document.getElementById('selected-images')
    selectedImages.innerHTML = ''
    fileInput = document.getElementById('images')
    let cloneCounter = 0
    fileInput.addEventListener('change', (e) => {
        files = e.target.files
        Array.from(files).forEach(file => {
            if (Images.length >= maxFileCount) {
                alert('You can only upload up to 8 images')
                return
            }
            Images.push(file)
            var template = document.getElementById('selected-file-template')
            var clone = template.content.cloneNode(true)
            let id = "selected-file-" + cloneCounter++
            clone.querySelector('div').id = id
            var fileNameElement = clone.querySelector('.file-name')
            fileNameElement.textContent = file.name
            var removeButton = clone.querySelector('.remove-button')
            removeButton.addEventListener('click', (e) => {
                e.preventDefault()
                Images = Images.filter(image => image != file)
                document.getElementById(id).remove()
            })
            selectedImages.appendChild(clone)
        })
        e.target.value = null
    })
    postForm = document.getElementById('post-form')
    postForm.addEventListener('submit', async (e) => {
        e.preventDefault()
        var formData = new FormData(e.target)
        if (Images.length <= 0) {
            alert('You must select at least one image')
            e.preventDefault()
            return
        }
        if (Images.length > maxFileCount) {
            alert('You can only upload up to 8 images')
            e.preventDefault()
            return
        }
        Images.forEach(image => {
            console.log('image :', image)
            if (image.size > maxFileSize) {
                alert('File size must be less than 5MB')
                e.preventDefault()
                return
            }
            formData.append('images', image)
        })
        console.log('formData :', Array.from(formData.entries()))
        await fetch('/api/v1/post/upload', {
            method: 'POST',
            body: formData
        }).then(res => {
            if (!res.ok) {
                alert(res.text)
            } else {
                Images = []
                selectedImages.innerHTML = ''
            }
        })

    })
    // postForm.addEventListener('htmx:configRequest', (e) => {
    //     console.log('before :', e.detail.body)
    //     var formData = new FormData(e.target)
    //     if (Images.length <= 0) {
    //         alert('You must select at least one image')
    //         e.preventDefault()
    //         return
    //     }
    //     if (Images.length > maxFileCount) {
    //         alert('You can only upload up to 8 images')
    //         e.preventDefault()
    //         return
    //     }
    //     Images.forEach(image => {
    //         console.log('image :', image)
    //         if (image.size > maxFileSize) {
    //             alert('File size must be less than 5MB')
    //             e.preventDefault()
    //             return
    //         }
    //         formData.append('images', image)
    //     })
    //     e.detail.body = formData
    //     console.log('after :', e.detail.body)
    // })


})