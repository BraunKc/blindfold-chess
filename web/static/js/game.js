class Game {
  constructor(api, stateView, promotionModal, notifier) {
    this.api = api
    this.stateView = stateView
    this.promotionModal = promotionModal
    this.notifier = notifier
  }

  async move(notation) {
    if (!notation) return

    try {
      const res = await this.api.isPromotion(notation)
      if (res.is_promotion) {
        const piece = await this.promotionModal.open(this.state.turn ?? 'w')
        if (!piece) return

        notation += piece
      }

      const state = await this.api.move(notation)
      this.state = state
      this.stateView.render(state)
      return true
    } catch (err) {
      this.notifier.showError(err.message)
      return false
    }
  }

  async init() {
    this.state = await this.api.state()
    this.stateView.render(this.state)
  }
}

class MoveEntry {
  constructor(game, input, cells, moveBtn) {
    this.game = game
    this.input = input
    this.cells = cells
    this.moveBtn = moveBtn

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

  async #onEnter(e) {
    e.preventDefault()

    const cell = e.target
    if (cell.value === '') return

    const notation = this.#getNotation()
    if (!notation) {
      const nextCell = this.#getNextCell(cell)
      if (nextCell) nextCell.focus()

      return
    }

    const ok = await this.game.move(notation)
    if (ok) {
      this.clear()
      this.cells[0].focus()
    }
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
    if (cell.value) {
      const nextCell = this.#getNextCell(cell)
      if (nextCell && !nextCell.value) nextCell.focus()
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

    this.moveBtn.addEventListener('click', () =>
      this.game.move(this.#getNotation()),
    )
  }
}

class StateView {
  constructor(turnEl, lastMovesEl) {
    this.turn = turnEl
    this.lastMoves = lastMovesEl
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
    if (!turn) return
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
    if (!state) return

    this.renderTurn(state.turn)

    if (this.lastMoves.children.length === 0) {
      this.renderLastMoves(state.last_moves)
    } else {
      this.renderLastMove(state.last_moves.at(-1))
    }
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
