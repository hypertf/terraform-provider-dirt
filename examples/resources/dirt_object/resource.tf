resource "dirt_bucket" "images" {
  name = "images_bucket"
}

resource "dirt_object" "logo" {
  bucket_id      = dirt_bucket.images.id
  path           = "assets/logo.png"
  content_base64 = filebase64("${path.module}/files/logo.png")
}
