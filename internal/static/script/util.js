document.addEventListener('alpine:init', () => {
    Alpine.store('postUploading', false)
})

var periodicIntersectUpdateObserver = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            entry.target.intervalID = setInterval(() => {
                htmx.trigger(entry.target, 'update')
            }, 5000)
        } else {
            clearInterval(entry.target.intervalID)
        }
    })
})

document.addEventListener('DOMContentLoaded', () => {
    const maxFileSize = 1024 * 1024 * 5
    const maxFileCount = 8
    var Images = []
    var selectedImages = document.getElementById('selected-images')
    selectedImages.innerHTML = ''
    fileInput = document.getElementById('images')
    let cloneCounter = 0
    fileInput.addEventListener('change', (e) => {
        files = e.target.files
        if (Images.length + files.length > maxFileCount) {
            alert('You can only upload up to 8 images')
            return
        }
        Array.from(files).forEach(file => {
            Images.push(file)
            var url = URL.createObjectURL(file)
            var template = document.getElementById('selected-file-template')
            var clone = template.content.cloneNode(true)
            let id = "selected-file-" + cloneCounter++
            clone.querySelector('div').id = id
            var imageElement = clone.querySelector('.file-image')
            imageElement.src = url
            imageElement.href = url
            imageElement.target = "_blank"
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
        Alpine.store('postUploading', true)
        var formData = new FormData(e.target)
        if (Images.length <= 0) {
            alert('You must select at least one image')
            e.preventDefault()
            Alpine.store('postUploading', false)
            return
        }
        if (Images.length > maxFileCount) {
            alert('You can only upload up to 8 images')
            e.preventDefault()
            Alpine.store('postUploading', false)
            return
        }
        Images.forEach(image => {
            if (image.size > maxFileSize) {
                alert('File size must be less than 5MB')
                e.preventDefault()
                Alpine.store('postUploading', false)
                return
            }
            formData.append('images', image)
        })
        await fetch('/api/v1/post/upload', {
            method: 'POST',
            body: formData
        }).then(res => {
            if (!res.ok) {
                res.text().then(txt => { alert(txt) })
            } else {
                Images = []
                selectedImages.innerHTML = ''
                e.target.reset()
                loader = document.getElementById("newer-post-loader")
                if (loader) {
                    htmx.trigger(loader, 'update')
                }
            }
            Alpine.store('postUploading', false)
        })

    })
})

function reactAndRefresh(e) {
    if (e.detail.xhr.status != 200) {
        return
    }
    var elem = e.target
    var parent = elem.closest('.react-section')
    var dislike = parent.querySelector('.dislike')
    var like = parent.querySelector('.like')

    if (!like || !dislike) {
        console.log("parent", parent)
        console.log('like or dislike not found')
        return
    }

    var likeurl = like.getAttribute('hx-post')
    var dislikeurl = dislike.getAttribute('hx-post')

    if (elem == like) {
        dislike.classList.remove('text-primary')
        if (likeurl.includes("value=1")) {
            likeurl = likeurl.replace('value=1', 'value=0')
        } else {
            likeurl = likeurl.replace('value=0', 'value=1')
        }
        if (dislikeurl.includes("value=0")) {
            dislikeurl = dislikeurl.replace('value=0', 'value=-1')
        }
    } else {
        like.classList.remove('text-primary')
        if (dislikeurl.includes("value=-1")) {
            dislikeurl = dislikeurl.replace('value=-1', 'value=0')
        } else {
            dislikeurl = dislikeurl.replace('value=0', 'value=-1')
        }
        if (likeurl.includes("value=0")) {
            likeurl = likeurl.replace('value=0', 'value=1')
        }
    }
    elem.classList.toggle('text-primary')

    like.setAttribute('hx-post', likeurl)
    dislike.setAttribute('hx-post', dislikeurl)

    htmx.process(like)
    htmx.process(dislike)

    var likeCount = parent.querySelector('.like-count')
    htmx.trigger(likeCount, 'update')
}
