import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';

class Square extends React.Component {
    constructor(){
	super();
    }
  render() {
    return (
	    <button className="square" onClick={() => this.props.onClick()}>
	    {this.props.value}
      </button>
    );
  }
}

class Board extends React.Component {
    constructor(){
	super();
	this.state = {turn: 0, squares: [null, null, null, null, null, null, null, null, null]};
    }
    renderSquare(i) {
	return <Square value={this.state.squares[i]} onClick={() => this.doClick(i)}/>;
    }
    getPlayer(){
	return this.state.turn % 2 == 0 ? 'x' : 'o';
    }
    doClick(i){
	if(this.state.squares[i] == null){
	    var s = this.state.squares.slice();
	    s[i] = this.getPlayer();
	    this.setState({turn: this.state.turn + 1, squares:s});
	}
    }
    
  render() {
      const status = 'Next player: ' + this.getPlayer();

    return (
      <div>
        <div className="status">{status}</div>
        <div className="board-row">
          {this.renderSquare(0)}
          {this.renderSquare(1)}
          {this.renderSquare(2)}
        </div>
        <div className="board-row">
          {this.renderSquare(3)}
          {this.renderSquare(4)}
          {this.renderSquare(5)}
        </div>
        <div className="board-row">
          {this.renderSquare(6)}
          {this.renderSquare(7)}
          {this.renderSquare(8)}
        </div>
      </div>
    );
  }
}

class Game extends React.Component {
  render() {
    return (
      <div className="game">
        <div className="game-board">
          <Board />
        </div>
        <div className="game-info">
          <div>{/* status */}</div>
          <ol>{/* TODO */}</ol>
        </div>
      </div>
    );
  }
}

// ========================================

ReactDOM.render(
  <Game />,
  document.getElementById('root')
);

