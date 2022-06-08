import React from 'react';
import PropTypes from 'prop-types';

const Root = ({ visible, close }) => {
  const [values, setValues] = React.useState({
    text: ''
  });

  if (!visible) {
    return null;
  }

  const style = getStyle();

  const handleChange = (name) => (event) => {
    const newValues = {
      ...values,
      [name]: event.target.value
    };
    setValues(newValues);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (values.text) {
      fetch('http://localhost:8065/plugins/com.example.mattermost-plugin-sample', {
        method: 'POST',
        body: JSON.stringify(values),
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        }
      }).then((res) => res.json()).then(console.log).catch(console.error);
    }
    setValues({ text: '' });
    close();
  };

  return (
    <div style={style.backdrop}>
      <div style={style.modal}>
        <div className='modal-content'>
          <div className='modal-header' style={{ backgroundColor: '#333' }}>
            <h1 className='modal-title' style={style.headertext}>Similar Words Search</h1>
          </div>
          <div className='modal-body'>
            <form onSubmit={handleSubmit}>
              <input type='text' style={style.searchinput} name='text' value={values.text} onChange={handleChange('text')} />
              <button type='button' className='btn btn-primary' style={style.searchbtn} onClick={handleSubmit}>
                <span>Search</span>
              </button>
            </form>
          </div>
          <div className='modal-footer'>
            <a style={{ textAlign: 'right', display: 'block' }} onClick={close}>Close</a>
          </div>
        </div>
      </div>
    </div>
  );
};

Root.propTypes = {
  visible: PropTypes.bool.isRequired,
  close: PropTypes.func.isRequired,
};

const getStyle = () => ({
  backdrop: {
    position: 'absolute',
    display: 'flex',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(0, 0, 0, 0.50)',
    zIndex: 2000,
    alignItems: 'baseline',
    justifyContent: 'center',
  },
  modal: {
    height: '150px',
    width: '600px',
    backgroundColor: 'white',
    marginTop: '100px',
  },
  headertext: {
    color: 'white',
    fontSize: '17px',
    lineHeight: '27px',
  },
  searchinput: {
    width: 200,
    height: '34px',
    padding: '6px 12px',
    borderRadius: '2px',
    border: '1px solid #ccc',
    color: '#555',
  },
  searchbtn: {
    marginLeft: '1em',
    marginBottom: '4px',
  }
});

export default Root;