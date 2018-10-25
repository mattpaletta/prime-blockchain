import * as React from 'react';
import './App.css';

class App extends React.Component<{}, { isView: boolean }> {
  constructor(props: any){
    super(props);
    this.state = { 
      isView: true,
    };
    this.viewCallback = this.viewCallback.bind(this);
    this.learnCallback = this.learnCallback.bind(this);
  }

  public viewCallback() {
    this.setState({ isView: true });
  }

  public learnCallback() {
    this.setState({ isView: false });
  }

  public render() {
    const viewElement = document.getElementById('view');
    const learnElement = document.getElementById('learn');

    if (viewElement && learnElement) {
      if (this.state.isView) {
        learnElement.style.textDecoration = 'none';
        viewElement.style.textDecoration = 'underline';
      } else {
        viewElement.style.textDecoration = 'none';
        learnElement.style.textDecoration = 'underline';
      }
    }

    return (
      <div className="heading-div">
        <div className='parent-div-1'>
          <div className='child-div'>prime blockchain</div>
        </div>

        <div className='parent-div-2'>
          <a id='view' className='child-div' onClick={ this.viewCallback }>view</a>
          <a id='learn' className='child-div' onClick={ this.learnCallback }>learn</a>
        </div>
      </div>
    );
  }
}

export default App;
