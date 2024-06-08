import React from 'react'
import styled from 'styled-components'
import { features } from '../lib/features';
import Game from '../lib/Game';
import Card from './Card';

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
      <Cards>
        {game.deck.map((card, index) => (
          <li key={index}>
            <Card card={card} optionsNumber={optionsNumber} />
          </li>
        ))}
      </Cards>
    </Wrapper>
  )
}

const Wrapper = styled.div`
  padding: 1rem 2rem;
`;

const Cards = styled.ul`
  display: flex;
  flex-wrap: wrap;
  list-style: none;
`;