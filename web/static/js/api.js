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
    })
  }

  async save() {
    return this.#request('/api/save', {
      method: 'POST',
    })
  }

  async pgn() {
    return this.#request('/api/pgn')
  }

  async turn() {
    return this.#request('/api/turn')
  }

  async lastMoves(lastMovesCount = 3) {
    return this.#request(`/api/last-moves?lmc=${lastMovesCount}`)
  }

  async outcome() {
    return this.#request('/api/outcome')
  }

  async isOver() {
    return this.#request('/api/is-over')
  }

  async state(lastMovesCount = 3) {
    return this.#request(`/api/state?lmc=${lastMovesCount}`)
  }

  async isPromotion(notation) {
    return this.#request(
      `/api/is-promotion?notation=${encodeURIComponent(notation)}`,
    )
  }
}
