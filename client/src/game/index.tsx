import React from 'react'
import styled from 'styled-components'
import { features } from '../lib/features';
import Game from '../lib/Game';

export default function GameWindow() {
  const gameFeatures = [
    features.shape,
    features.color,
    features.number,
    features.shading,
    // features.rotation,
  ]
  const optionsNumber = 3

  const game = new Game({ features: gameFeatures, optionsNumber})
  console.log(game)
  return (
    <Wrapper>
      Game
      <ul>
        {game.deck.map((card, index) => (
          <li key={index}>
            {game.features.map((feature) => (
              <span key={feature}>
                {feature}: {String(card[feature])}
              </span>
            ))}
            <br/>
          </li>
        ))}
      </ul>
    </Wrapper>
  )
}

const Wrapper = styled.div`

`;