class Notifier {
  constructor(el, duration = 4000) {
    this.el = el
    this.duration = duration
    this.timer = null
  }

  #show(msg, modifier) {
    this.el.textContent = msg
    this.el.classList.remove('hide', 'error', 'success')
    if (modifier) this.el.classList.add(modifier)

    clearTimeout(this.timer)
    this.timer = setTimeout(() => {
      this.el.classList.add('hide')
    }, this.duration)
  }

  showError(msg) {
    this.#show(msg, 'error')
  }

  showSuccess(msg) {
    this.#show(msg, 'success')
  }

  showMessage(msg) {
    this.#show(msg)
  }
}
