class API {
  async #request(url, options = {}) {
    const res = await fetch(url, options)
    const text = await res.text()
    const data = text ? JSON.parse(text) : {}
    if (!res.ok) throw new Error(data.error || `HTTP ${res.status}`)

    return data
  }

  async move(notation) {
    if (!notation) return

    return this.#request('/api/move', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        notation: notation,
      }),
    })
  }

  async reset() {
    return this.#request('/api/reset', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    })
  }

  async save() {
    return this.#request('/api/save', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    })
  }

  async pgn() {
    return this.#request('/api/pgn')
  }

  async turn() {
    return this.#request('/api/turn')
  }

  async lastMoves(count = 3) {
    return this.#request(`/api/last-moves?lmc=${count}`)
  }

  async outcome() {
    return this.#request('/api/outcome')
  }

  async isOver() {
    return this.#request('/api/is-over')
  }

  async state() {
    return this.#request('/api/state')
  }

  async isPromotion(notation) {
    return this.#request(`/api/is-promotion?notation=${notation}`)
  }
}

class MoveEntry {
  constructor(api, stateView, promotionModal, input, cells, moveBtn, notifier) {
    this.api = api
    this.state = stateView
    this.promotionModal = promotionModal
    this.input = input
    this.cells = cells
    this.moveBtn = moveBtn
    this.notifier = notifier
    this.lastState = null

    this.#addListeners()
  }

  clear() {
    this.cells.forEach((el) => {
      el.value = ''
    })
  }

  #getNextCell(cell) {
    const cellIndex = Number(cell.dataset.index)
    return cellIndex < this.cells.length - 1
      ? this.cells[cellIndex + 1]
      : undefined
  }

  #getPrevCell(cell) {
    const cellIndex = Number(cell.dataset.index)
    return cellIndex > 0 ? this.cells[cellIndex - 1] : undefined
  }

  #getNotation() {
    let notation = ''
    this.cells.forEach((el) => {
      notation += el.value.toLowerCase()
    })

    return notation.length >= 4 ? notation : null
  }

  async #submitMove() {
    let notation = this.#getNotation()
    if (!notation) return

    try {
      const isPromotion = await this.api.isPromotion(notation)
      if (isPromotion.is_promotion) {
        const piece = await this.promotionModal.open(
          this.lastState?.turn ?? 'w',
        )
        if (!piece) return
        notation += piece
      }

      const state = await this.api.move(notation)
      this.lastState = state
      this.clear()
      this.cells[0].focus()
      this.state.render(state)
    } catch (err) {
      this.notifier.showError(err.message)
    }
  }

  #validateInput(e) {
    const cell = e.target
    const cellIndex = Number(cell.dataset.index)
    let value = cell.value.slice(0, 1)

    const isLetterCell = cellIndex % 2 === 0
    const valid = isLetterCell ? /^[a-h]$/i.test(value) : /^[1-8]$/.test(value)

    if (value && !valid) {
      cell.value = ''
      return
    }

    cell.value = value.toLowerCase()
  }

  #onEnter(e) {
    e.preventDefault()

    const cell = e.target
    if (cell.value === '') return

    this.#submitMove()

    const nextCell = this.#getNextCell(cell)
    if (nextCell) nextCell.focus()
  }

  #onBackspace(e) {
    e.preventDefault()

    const cell = e.target
    if (cell.value !== '') {
      cell.value = ''
      return
    }

    const prevCell = this.#getPrevCell(cell)
    if (!prevCell) return

    prevCell.focus()
    prevCell.value = ''
  }

  #onArrows(e) {
    e.preventDefault()

    const cellIndex = Number(e.target.dataset.index)
    const next = e.key === 'ArrowLeft' ? cellIndex - 1 : cellIndex + 1

    if (next >= 0 && next < this.cells.length) {
      this.cells[next].focus()
    }
  }

  #onFocus(e) {
    const len = e.target.value.length
    e.target.setSelectionRange(len, len)
  }

  #onInput(e) {
    this.#validateInput(e)

    const cell = e.target
    if (cell.value !== '') {
      const nextCell = this.#getNextCell(cell)
      if (nextCell && nextCell.value === '') nextCell.focus()
    }
  }

  #onKeydown(e) {
    if (e.key === 'Enter') this.#onEnter(e)
    else if (e.key === 'Backspace') this.#onBackspace(e)
    else if (e.key === 'ArrowLeft' || e.key === 'ArrowRight') this.#onArrows(e)
  }

  #addListeners() {
    this.input.addEventListener('focusin', (e) => this.#onFocus(e))
    this.input.addEventListener('input', (e) => this.#onInput(e))
    this.input.addEventListener('keydown', (e) => this.#onKeydown(e))

    this.moveBtn.addEventListener('click', () => this.#submitMove())
  }
}

class StateView {
  constructor(turn, lastMoves) {
    this.turn = turn
    this.lastMoves = lastMoves
  }

  #getActiveMoves() {
    return [...this.lastMoves.children].filter(
      (el) => !el.classList.contains('hide'),
    )
  }

  #addMove(notation) {
    const li = document.createElement('li')
    li.innerText = notation
    li.classList.add('hide')
    this.lastMoves.prepend(li)

    li.offsetHeight
    li.classList.remove('hide')
  }

  #removeMove() {
    const items = this.#getActiveMoves()
    const lastEl = items[items.length - 1]
    if (!lastEl) return

    lastEl.classList.add('hide')
    lastEl.addEventListener('transitionend', () => lastEl.remove(), {
      once: true,
    })
  }

  renderTurn(turn) {
    this.turn.innerText = turn === 'w' ? 'ход белых' : 'ход черных'
  }

  renderLastMove(notation, count = 3) {
    if (!notation) return

    this.#addMove(notation)

    const excess = this.#getActiveMoves().length - count
    for (let i = 0; i < excess; i++) {
      this.#removeMove()
    }
  }

  renderLastMoves(lastMoves, count = 3) {
    if (!lastMoves) return

    lastMoves.forEach((el) => {
      this.#addMove(el)
    })

    const excess = this.#getActiveMoves().length - count
    for (let i = 0; i < excess; i++) {
      this.#removeMove()
    }
  }

  render(state) {
    this.renderTurn(state.turn)

    if (this.lastMoves.children.length === 0) {
      this.renderLastMoves(state.last_moves)
    } else {
      this.renderLastMove(state.last_moves.at(-1))
    }
  }
}

class Notifier {
  constructor(el, duration = 4000) {
    this.el = el
    this.duration = duration
    this.timer = null
  }

  #show(msg, modifier) {
    this.el.textContent = msg
    this.el.classList.remove('hide', 'error', 'success')
    this.el.classList.add(modifier)

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

class PromotionModal {
  #res = null

  constructor(modal, backdrop) {
    this.modal = modal
    this.backdrop = backdrop

    this.#addListeners()
  }

  open(color) {
    this.modal.classList.remove('white', 'black')
    this.modal.classList.add(color === 'w' ? 'white' : 'black')
    this.modal.classList.remove('hide')
    this.backdrop.classList.remove('hide')

    return new Promise((res) => {
      this.#res = res
    })
  }

  #close(piece) {
    this.modal.classList.add('hide')
    this.backdrop.classList.add('hide')

    if (this.#res) {
      this.#res(piece)
      this.#res = null
    }
  }

  #addListeners() {
    this.modal.addEventListener('click', (e) => {
      const btn = e.target.closest('[data-piece]')
      if (!btn) return
      this.#close(btn.dataset.piece)
    })

    this.backdrop.addEventListener('click', () => this.#close())
  }
}

const api = new API()

const notificationEl = document.getElementById('notification')
const notifier = new Notifier(notificationEl)

const turnEl = document.getElementById('turn')
const lastMovesEl = document.getElementById('lastMoves')
const stateView = new StateView(turnEl, lastMovesEl)

const promotionModalEl = document.getElementById('promotionModal')
const promotionBackdropEl = document.getElementById('promotionBackdrop')
const promotionModal = new PromotionModal(promotionModalEl, promotionBackdropEl)

const inputEl = document.getElementById('moveInput')
const cellsEls = document.querySelectorAll('.move-cell')
const moveBtnEl = document.getElementById('moveBtn')
const moveEntry = new MoveEntry(
  api,
  stateView,
  promotionModal,
  inputEl,
  cellsEls,
  moveBtnEl,
  notifier,
)

document.addEventListener('DOMContentLoaded', async () => {
  try {
    const state = await api.state()
    stateView.render(state)
  } catch (err) {
    notifier.showError(err.message)
  }
})

// document.addEventListener('click', () => {
//   document.getElementById('promotionModal').classList.toggle('hide')
//   document.getElementById('promotionModal').classList.toggle('white')
//   document.getElementById('promotionBackdrop').classList.toggle('hide')
// })
