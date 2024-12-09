import React from 'react';
import { render, screen } from '@testing-library/react';
import App from './App';

test('renders page correctly', () => {
  render(<App />);
  const searchheader = screen.getByText(/Search Number/i);
  expect(searchheader).toBeInTheDocument();
});
