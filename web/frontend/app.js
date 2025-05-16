console.log('Loading app.js...');

if (typeof React === 'undefined') {
    console.error('React is not loaded!');
    document.getElementById('root').innerHTML = '<div class="alert alert-danger m-5">Error: React is not loaded</div>';
    throw new Error('React is not loaded');
}

if (typeof ReactDOM === 'undefined') {
    console.error('ReactDOM is not loaded!');
    document.getElementById('root').innerHTML = '<div class="alert alert-danger m-5">Error: ReactDOM is not loaded</div>';
    throw new Error('ReactDOM is not loaded');
}

console.log('React version:', React.version);
console.log('ReactDOM version:', ReactDOM.version);

window.onerror = function(message, source, lineno, colno, error) {
    console.error('Global error:', { message, source, lineno, colno, error });
    return false;
};

const AppContext = React.createContext(null);

function parseNumber(value, type = 'float', base = null) {
    value = value.trim().toLowerCase().replace(/\s+/g, '');
    
    if (!value) {
        switch (type) {
            case 'complex':
                return { type: 'complex', value: { real: 0, imag: 0 } };
            case 'gf':
                return { type: 'gf', value: 0 };
            default:
                return { type: 'float', value: 0 };
        }
    }

    if (type === 'gf') {
        if (!base) {
            console.error('Base is required for GF numbers');
            return { type: 'gf', value: 0 };
        }
        
        if (!/^-?\d+$/.test(value)) {
            return { type: 'gf', value: 0 };
        }
        
        let num = parseInt(value);
        
        if (isNaN(num)) {
            return { type: 'gf', value: 0 };
        }
        
        if (num < 0) {
            num = base - ((-num) % base);
        } else {
            num = num % base;
        }
        
        return { 
            type: 'gf', 
            value: num
        };
    }
    
    if (type === 'complex') {
        if (!value.includes('i')) {
            const num = parseFloat(value) || 0;
            return { 
                type: 'complex', 
                value: { real: num, imag: 0 } 
            };
        }

        if (value === 'i') return { type: 'complex', value: { real: 0, imag: 1 } };
        if (value === '-i') return { type: 'complex', value: { real: 0, imag: -1 } };
        if (value === '+i') return { type: 'complex', value: { real: 0, imag: 1 } };

        value = value
            .replace(/^i/, '1i')
            .replace(/([+-])i(?!\d)/, '$11i');

        let real = 0;
        let imag = 0;

        const matches = value.match(/([+-]?\d*\.?\d*)(i)?/g).filter(Boolean);
        
        for (const match of matches) {
            if (match.includes('i')) {
                const num = match.replace('i', '');
                if (num === '' || num === '+') imag = 1;
                else if (num === '-') imag = -1;
                else imag = parseFloat(num);
            } else {
                real = parseFloat(match);
            }
        }

        return { 
            type: 'complex', 
            value: { 
                real: isNaN(real) ? 0 : real, 
                imag: isNaN(imag) ? 0 : imag 
            } 
        };
    }

    const num = parseFloat(value);
    if (!isNaN(num)) {
        return { type: 'float', value: num };
    }
    
    return null;
}

function formatNumber(number) {
    if (!number) return '0';
    
    switch (number.type) {
        case 'complex': {
            const real = number.value.real || 0;
            const imag = number.value.imag || 0;
            
            const roundedReal = Number(real.toFixed(2));
            const roundedImag = Number(imag.toFixed(2));
            
            if (roundedReal === 0 && roundedImag === 0) return '0';
            if (roundedReal === 0) return roundedImag === 1 ? 'i' : roundedImag === -1 ? '-i' : `${roundedImag}i`;
            if (roundedImag === 0) return `${roundedReal}`;
            
            const imagPart = roundedImag === 1 ? 'i' : roundedImag === -1 ? '-i' : `${Math.abs(roundedImag)}i`;
            return `${roundedReal}${roundedImag < 0 ? '-' : '+'}${imagPart}`;
        }
        case 'gf':
            return number.value.toString();
        case 'float':
            return Number(number.value.toFixed(2)).toString();
        default:
            return number.value.toString();
    }
}

function formatNumberForServer(number) {
    if (!number) return '0';
    
    switch (number.type) {
        case 'complex': {
            const real = number.value.real || 0;
            const imag = number.value.imag || 0;
            
            if (real === 0 && imag === 0) return '0';
            if (real === 0) return imag === 1 ? 'i' : imag === -1 ? '-i' : `${imag}i`;
            if (imag === 0) return `${real}`;
            
            const imagPart = imag === 1 ? 'i' : imag === -1 ? '-i' : `${Math.abs(imag)}i`;
            return `${real}${imag < 0 ? '-' : '+'}${imagPart}`;
        }
        case 'gf':
            return number.value.toString();
        case 'float':
            return number.value.toString();
        default:
            return number.value.toString();
    }
}

function NumberInput({ value, onChange, readOnly, onKeyDown }) {
    const [inputValue, setInputValue] = React.useState(formatNumber(value));
    const [isFocused, setIsFocused] = React.useState(false);
    const [previousValue, setPreviousValue] = React.useState(null);
    const inputRef = React.useRef(null);

    const { numberType, fieldBase } = React.useContext(AppContext) || { numberType: 'float', fieldBase: 2 };
    console.log('NumberInput context:', { numberType, fieldBase });

    React.useEffect(() => {
        if (!isFocused) {
            const formatted = formatNumber(value);
            console.log('NumberInput effect:', { value, formatted });
            setInputValue(formatted);
        }
    }, [value, isFocused]);

    const handleFocus = () => {
        if (readOnly) return;
        setIsFocused(true);
        setPreviousValue(inputValue);
        setInputValue('');
    };

    const handleBlur = () => {
        if (readOnly) return;
        setIsFocused(false);
        if (inputValue === '') {
            setInputValue(previousValue);
            const parsedNumber = parseNumber(previousValue, numberType, fieldBase);
            if (parsedNumber) {
                onChange(parsedNumber);
            }
        }
    };

    const handleChange = (e) => {
        if (readOnly) return;
        const newValue = e.target.value;
        console.log('NumberInput handleChange:', { newValue, numberType, fieldBase });
        setInputValue(newValue);
        
        const parsedNumber = parseNumber(newValue, numberType, fieldBase);
        console.log('NumberInput parsed:', parsedNumber);
        if (parsedNumber) {
            onChange(parsedNumber);
        }
    };

    return (
        <input
            ref={inputRef}
            type="text"
            className={`matrix-input form-control ${readOnly ? 'readonly-result' : ''}`}
            value={inputValue}
            onChange={handleChange}
            onFocus={handleFocus}
            onBlur={handleBlur}
            onKeyDown={onKeyDown}
            readOnly={readOnly}
            placeholder="0"
            style={{
                ...readOnly ? { cursor: 'default', backgroundColor: '#f8f9fa', minWidth: '120px' } : { minWidth: '120px' },
                WebkitAppearance: 'none'
            }}
        />
    );
}

const numberInputStyle = {
    '::-webkit-outer-spin-button': {
        '-webkit-appearance': 'none',
        margin: 0,
    },
    '::-webkit-inner-spin-button': {
        '-webkit-appearance': 'none',
        margin: 0,
    },
    MozAppearance: 'textfield',
    '::-webkit-search-cancel-button': {
        '-webkit-appearance': 'none',
    }
};

function Matrix({ rows, cols, values, onChange, readOnly }) {
    const handleKeyDown = (i, j) => (event) => {
        if (readOnly) return;

        const moveMap = {
            'ArrowUp': [-1, 0],
            'ArrowDown': [1, 0],
            'ArrowLeft': [0, -1],
            'ArrowRight': [0, 1],
            'Tab': [0, 1],
            'Enter': [1, 0]
        };

        if (event.key in moveMap) {
            event.preventDefault();
            const [di, dj] = moveMap[event.key];
            let newI = i + di;
            let newJ = j + dj;

            if (event.key === 'Tab' && event.shiftKey) {
                newJ = j - 1;
            }

            if (event.key === 'Tab') {
                if (newJ >= cols) {
                    newJ = 0;
                    newI++;
                } else if (newJ < 0) {
                    newJ = cols - 1;
                    newI--;
                }
            }

            if (event.key === 'Enter') {
                if (newI >= rows) {
                    newI = 0;
                }
            }

            if (newI >= 0 && newI < rows && newJ >= 0 && newJ < cols) {
                const nextInput = document.querySelector(
                    `.matrix-row:nth-child(${newI + 1}) .matrix-input:nth-child(${newJ + 1})`
                );
                if (nextInput) {
                    nextInput.focus();
                }
            }
        }
    };

    return (
        <div className="matrix-container">
            <div className="d-flex">
                <div className="matrix-bracket">[</div>
                <div>
                    {Array(rows).fill().map((_, i) => (
                        <div key={i} className="matrix-row">
                            {Array(cols).fill().map((_, j) => (
                                <NumberInput
                                    key={j}
                                    value={values[i][j]}
                                    onChange={(value) => onChange(i, j, value)}
                                    onKeyDown={handleKeyDown(i, j)}
                                    readOnly={readOnly}
                                />
                            ))}
                        </div>
                    ))}
                </div>
                <div className="matrix-bracket">]</div>
            </div>
        </div>
    );
}

function Vector({ size, values, onChange, readOnly }) {
    const handleKeyDown = (i) => (event) => {
        if (readOnly) return;

        if (event.key === 'ArrowUp' && i > 0) {
            event.preventDefault();
            const prevInput = document.querySelector(
                `.vector-row:nth-child(${i}) .matrix-input`
            );
            if (prevInput) prevInput.focus();
        } else if (event.key === 'ArrowDown' && i < size - 1) {
            event.preventDefault();
            const nextInput = document.querySelector(
                `.vector-row:nth-child(${i + 2}) .matrix-input`
            );
            if (nextInput) nextInput.focus();
        } else if (event.key === 'Enter') {
            event.preventDefault();
            const nextIndex = (i + 1) % size;
            const nextInput = document.querySelector(
                `.vector-row:nth-child(${nextIndex + 1}) .matrix-input`
            );
            if (nextInput) nextInput.focus();
        }
    };

    return (
        <div className="matrix-container">
            <div className="d-flex">
                <div className="matrix-bracket">[</div>
                <div>
                    {Array(size).fill().map((_, i) => (
                        <div key={i} className="matrix-row vector-row">
                            <NumberInput
                                value={values[i]}
                                onChange={(value) => onChange(i, value)}
                                onKeyDown={handleKeyDown(i)}
                                readOnly={readOnly}
                            />
                        </div>
                    ))}
                </div>
                <div className="matrix-bracket">]</div>
            </div>
        </div>
    );
}

function LoadingSpinner() {
    return (
        <div className="d-flex justify-content-center align-items-center" style={{ minHeight: '200px' }}>
            <div className="spinner-border text-primary" role="status">
                <span className="visually-hidden">Загрузка...</span>
            </div>
        </div>
    );
}

function OperationSelect({ value, onChange }) {
    const [isOpen, setIsOpen] = React.useState(false);
    const selectRef = React.useRef(null);

    const operations = [
        { value: 'solve', label: 'Решить систему уравнений', icon: '🧮' },
        { value: 'determinant', label: 'Вычислить определитель', icon: '📊' },
        { value: 'rank', label: 'Вычислить ранг', icon: '📈' },
        { value: 'inverse', label: 'Найти обратную матрицу', icon: '🔄' }
    ];

    const selectedOperation = operations.find(op => op.value === value);

    React.useEffect(() => {
        const handleClickOutside = (event) => {
            if (selectRef.current && !selectRef.current.contains(event.target)) {
                setIsOpen(false);
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    return (
        <div className="custom-select" ref={selectRef}>
            <div 
                className={`select-selected ${isOpen ? 'select-arrow-active' : ''}`}
                onClick={() => setIsOpen(!isOpen)}
            >
                {selectedOperation.icon} {selectedOperation.label}
            </div>
            {isOpen && (
                <div className="select-items">
                    {operations.map(op => (
                        <div 
                            key={op.value}
                            className={`select-item ${op.value === value ? 'same-as-selected' : ''}`}
                            onClick={() => {
                                onChange(op.value);
                                setIsOpen(false);
                            }}
                        >
                            {op.icon} {op.label}
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}

function SizeInput({ value, onChange, placeholder }) {
    const [inputValue, setInputValue] = React.useState(value.toString());
    const [isFocused, setIsFocused] = React.useState(false);
    const [previousValue, setPreviousValue] = React.useState(null);

    React.useEffect(() => {
        if (!isFocused) {
            setInputValue(value.toString());
        }
    }, [value, isFocused]);

    const handleFocus = () => {
        setIsFocused(true);
        setPreviousValue(inputValue);
        setInputValue('');
    };

    const handleBlur = () => {
        setIsFocused(false);
        if (inputValue === '') {
            setInputValue(previousValue);
            onChange(parseInt(previousValue) || 1);
        }
    };

    const handleChange = (e) => {
        const newValue = e.target.value;
        setInputValue(newValue);
        
        const numValue = parseInt(newValue);
        if (!isNaN(numValue) && numValue >= 1 && numValue <= 5) {
            onChange(numValue);
        }
    };

    return (
        <input
            type="number"
            className="form-control"
            value={inputValue}
            onChange={handleChange}
            onFocus={handleFocus}
            onBlur={handleBlur}
            min="1"
            max="5"
            placeholder={placeholder}
        />
    );
}

function isPrime(num) {
    if (num <= 1) return false;
    if (num === 2) return true;
    if (num % 2 === 0) return false;
    
    const sqrt = Math.sqrt(num);
    for (let i = 3; i <= sqrt; i += 2) {
        if (num % i === 0) return false;
    }
    return true;
}

function App() {
    console.log('App component rendering');
    
    const [rows, setRows] = React.useState(2);
    const [cols, setCols] = React.useState(2);
    const [matrixA, setMatrixA] = React.useState(Array(2).fill().map(() => Array(2).fill({ type: 'float', value: 0 })));
    const [vectorB, setVectorB] = React.useState(Array(2).fill({ type: 'float', value: 0 }));
    const [operation, setOperation] = React.useState('solve');
    const [result, setResult] = React.useState(null);
    const [error, setError] = React.useState(null);
    const [loading, setLoading] = React.useState(false);
    const [numberType, setNumberType] = React.useState('float');
    const [fieldBase, setFieldBase] = React.useState(2);
    const [fieldBaseInput, setFieldBaseInput] = React.useState('2');
    const [fieldBaseError, setFieldBaseError] = React.useState(null);

    React.useEffect(() => {
        if (typeof particlesJS !== 'undefined') {
            particlesJS('particles-js', {
                interactivity: {
                    detect_on: 'window',
                    events: {
                        onhover: {
                            enable: true,
                            mode: 'repulse'
                        },
                        onclick: {
                            enable: true,
                            mode: 'push'
                        },
                        resize: true
                    },
                    modes: {
                        repulse: {
                            distance: 100,
                            duration: 2.5,
                            speed: 10,
                            factor: 200,
                            maxSpeed: 80,
                            easing: 'cubic-bezier(0.25, 0.46, 0.45, 0.94)',
                            particles_momentum: true,
                            momentum_decay: 0.05
                        },
                        push: {
                            particles_nb: 1
                        }
                    }
                },
                particles: {
                    number: {
                        value: 180,
                        density: {
                            enable: true,
                            value_area: 1000
                        }
                    },
                    color: {
                        value: '#0d6efd'
                    },
                    shape: {
                        type: 'circle'
                    },
                    opacity: {
                        value: 0.5,
                        random: true,
                        anim: {
                            enable: false
                        }
                    },
                    size: {
                        value: 3,
                        random: true,
                        anim: {
                            enable: false
                        }
                    },
                    line_linked: {
                        enable: true,
                        distance: 150,
                        color: '#0d6efd',
                        opacity: 0.4,
                        width: 1
                    },
                    move: {
                        enable: true,
                        speed: 4,
                        direction: 'none',
                        random: true,
                        straight: false,
                        out_mode: 'bounce',
                        bounce: true,
                        attract: {
                            enable: true,
                            rotateX: 1200,
                            rotateY: 1800
                        },
                        decay: {
                            enable: true,
                            speed: 0.05,
                            min_speed: 0.5
                        },
                        gravity: {
                            enable: true,
                            acceleration: 9.8,
                            max_speed: 50
                        },
                        friction: {
                            enable: true,
                            value: 0.02
                        }
                    }
                },
                retina_detect: true,
                fps_limit: 60
            });
        }
    }, []);

    const needsSquareMatrix = ['solve', 'determinant', 'inverse'].includes(operation);

    const numberTypes = [
        { value: 'float', label: 'Действительные числа', icon: '🔢' },
        { value: 'complex', label: 'Комплексные числа', icon: '💫' },
        { value: 'gf', label: 'Конечное поле GF', icon: '🔄' }
    ];

    const fieldBases = [2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31];

    const convertToType = (num) => {
        if (typeof num === 'object' && num !== null) {
            switch (numberType) {
                case 'complex':
                    if (num.type === 'complex') return num;
                    return { 
                        type: 'complex', 
                        value: { 
                            real: num.type === 'float' ? num.value : 0, 
                            imag: 0 
                        } 
                    };
                case 'gf':
                    const val = num.type === 'complex' ? num.value.real : num.value;
                    return { type: 'gf', value: ((val % fieldBase) + fieldBase) % fieldBase };
                default:
                    return { type: 'float', value: num.type === 'complex' ? num.value.real : num.value };
            }
        }

        switch (numberType) {
            case 'complex':
                return { type: 'complex', value: { real: Number(num) || 0, imag: 0 } };
            case 'gf':
                return { type: 'gf', value: ((Number(num) % fieldBase) + fieldBase) % fieldBase };
            default:
                return { type: 'float', value: Number(num) || 0 };
        }
    };

    const validateSize = (size) => {
        if (isNaN(size) || size < 1) return 1;
        if (size > 5) return 5;
        return size;
    };

    const handleSizeChange = (newSize) => {
        const validSize = validateSize(newSize);
        if (validSize !== rows) {
            setRows(validSize);
            if (needsSquareMatrix) {
                setCols(validSize);
                setMatrixA(Array(validSize).fill().map(() => 
                    Array(validSize).fill(convertToType(0))
                ));
            } else {
                setMatrixA(Array(validSize).fill().map(() => 
                    Array(cols).fill(convertToType(0))
                ));
            }
            if (operation === 'solve') {
                setVectorB(Array(validSize).fill(convertToType(0)));
            }
            setResult(null);
            setError(null);
        }
    };

    const handleColsChange = (newCols) => {
        const validCols = validateSize(newCols);
        if (validCols !== cols) {
            setCols(validCols);
            setMatrixA(matrixA.map(row => {
                const newRow = [...row];
                while (newRow.length < validCols) {
                    newRow.push({ type: 'float', value: 0 });
                }
                return newRow.slice(0, validCols);
            }));
            setResult(null);
            setError(null);
        }
    };

    const handleMatrixChange = (i, j, value) => {
        const newMatrix = [...matrixA];
        newMatrix[i][j] = value;
        setMatrixA(newMatrix);
    };

    const handleVectorChange = (i, value) => {
        const newVector = [...vectorB];
        newVector[i] = value;
        setVectorB(newVector);
    };

    const handleNumberTypeChange = (newType) => {
        setNumberType(newType);
        const oldType = numberType;
        setMatrixA(matrixA.map(row => 
            row.map(num => {
                let currentValue;
                if (num.type === 'complex') {
                    currentValue = Math.round(num.value.real + (num.value.imag !== 0 ? num.value.imag : 0));
                } else {
                    currentValue = Math.round(num.value);
                }
                
                switch (newType) {
                    case 'complex':
                        return { 
                            type: 'complex', 
                            value: { 
                                real: Number(currentValue) || 0, 
                                imag: 0 
                            } 
                        };
                    case 'gf':
                        let fieldValue = currentValue;
                        if (fieldValue < 0) {
                            fieldValue = fieldBase - ((-fieldValue) % fieldBase);
                        } else {
                            fieldValue = fieldValue % fieldBase;
                        }
                        return { 
                            type: 'gf', 
                            value: fieldValue
                        };
                    default:
                        return { 
                            type: 'float', 
                            value: Number(currentValue) || 0 
                        };
                }
            })
        ));
        
        if (operation === 'solve') {
            setVectorB(vectorB.map(num => {
                let currentValue;
                if (num.type === 'complex') {
                    currentValue = Math.round(num.value.real + (num.value.imag !== 0 ? num.value.imag : 0));
                } else {
                    currentValue = Math.round(num.value);
                }
                
                switch (newType) {
                    case 'complex':
                        return { 
                            type: 'complex', 
                            value: { 
                                real: Number(currentValue) || 0, 
                                imag: 0 
                            } 
                        };
                    case 'gf':
                        let fieldValue = currentValue;
                        if (fieldValue < 0) {
                            fieldValue = fieldBase - ((-fieldValue) % fieldBase);
                        } else {
                            fieldValue = fieldValue % fieldBase;
                        }
                        return { 
                            type: 'gf', 
                            value: fieldValue
                        };
                    default:
                        return { 
                            type: 'float', 
                            value: Number(currentValue) || 0 
                        };
                }
            }));
        }
        setResult(null);
    };

    const handleFieldBaseChange = (e) => {
        const inputValue = e.target.value;
        setFieldBaseInput(inputValue);

        if (inputValue === '') {
            setFieldBaseError("Введите число");
            return;
        }

        const value = parseInt(inputValue);
        
        if (isNaN(value)) {
            setFieldBaseError("Введите целое число");
            return;
        }
        
        if (value <= 1) {
            setFieldBaseError("Основание поля должно быть больше 1");
            return;
        }
        
        if (!isPrime(value)) {
            setFieldBaseError("Основание поля должно быть простым числом");
            return;
        }
        
        setFieldBaseError(null);
        setFieldBase(value);
    };

    const handleSubmit = async () => {
        setLoading(true);
        setError(null);
        setResult(null);

        try {
            const matrixForServer = matrixA.map(row => 
                row.map(num => {
                    switch (numberType) {
                        case 'complex':
                            const real = num.value.real || 0;
                            const imag = num.value.imag || 0;
                            if (imag === 0) return real.toString();
                            if (imag > 0) return `${real}+${imag}i`;
                            return `${real}${imag}i`;
                        case 'gf':
                            return num.value.toString();
                        default:
                            return num.value.toString();
                    }
                })
            );

            const vectorForServer = operation === 'solve' ? 
                vectorB.map(num => {
                    switch (numberType) {
                        case 'complex':
                            const real = num.value.real || 0;
                            const imag = num.value.imag || 0;
                            if (imag === 0) return real.toString();
                            if (imag > 0) return `${real}+${imag}i`;
                            return `${real}${imag}i`;
                        case 'gf':
                            return num.value.toString();
                        default:
                            return num.value.toString();
                    }
                }) : 
                undefined;

            let requestData;
            if (operation === 'solve') {
                requestData = {
                    matrix: {
                        type: numberType === 'complex' ? 'complex' : 
                              numberType === 'gf' ? 'gf' : 'float64',
                        rows: rows,
                        cols: cols,
                        data: matrixForServer
                    }
                };
                requestData.vector = vectorForServer;
            } else {
                requestData = {
                    type: numberType === 'complex' ? 'complex' : 
                          numberType === 'gf' ? 'gf' : 'float64',
                    rows: rows,
                    cols: cols,
                    data: matrixForServer
                };
            }

            if (numberType === 'gf') {
                if (operation === 'solve') {
                    requestData.matrix.modP = fieldBase;
                } else {
                    requestData.modP = fieldBase;
                }
            }

            console.log('Sending request:', JSON.stringify(requestData, null, 2));

            const response = await fetch(`/api/${operation}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestData),
            });

            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || 'Произошла ошибка при вычислении');
            }

            const processValue = (val) => {
                switch (numberType) {
                    case 'complex':
                        if (typeof val === 'string') {
                            const matches = val.match(/([+-]?\d*\.?\d*)(i)?/g).filter(Boolean);
                            let real = 0;
                            let imag = 0;
                            
                            for (const match of matches) {
                                if (match.includes('i')) {
                                    const num = match.replace('i', '');
                                    if (num === '' || num === '+') imag = 1;
                                    else if (num === '-') imag = -1;
                                    else imag = parseFloat(num);
                                } else {
                                    real = parseFloat(match);
                                }
                            }
                            
                            return {
                                type: 'complex',
                                value: {
                                    real: isNaN(real) ? 0 : real,
                                    imag: isNaN(imag) ? 0 : imag
                                }
                            };
                        }
                        return { type: 'complex', value: { real: 0, imag: 0 } };
                    
                    case 'gf':
                        const gfVal = parseInt(val);
                        return {
                            type: 'gf',
                            value: ((gfVal % fieldBase) + fieldBase) % fieldBase
                        };
                    
                    default:
                        const num = parseFloat(val);
                        return {
                            type: 'float',
                            value: isNaN(num) ? 0 : num
                        };
                }
            };

            if (Array.isArray(data.result)) {
                if (Array.isArray(data.result[0])) {
                    setResult(data.result.map(row => row.map(processValue)));
                } else {
                    setResult(data.result.map(processValue));
                }
            } else if (data.value !== undefined) {
                setResult(processValue(data.value));
            } else {
                setResult(processValue(data.result));
            }
        } catch (err) {
            console.error('Error:', err);
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    const getOperationTitle = () => {
        switch (operation) {
            case 'solve': return 'Решение системы уравнений';
            case 'determinant': return 'Вычисление определителя';
            case 'rank': return 'Вычисление ранга';
            case 'inverse': return 'Обратная матрица';
            default: return '';
        }
    };

    const getOperationIcon = () => {
        switch (operation) {
            case 'solve': return 'bi-calculator-fill';
            case 'determinant': return 'bi-braces';
            case 'rank': return 'bi-bar-chart-fill';
            case 'inverse': return 'bi-arrow-repeat';
            default: return 'bi-calculator';
        }
    };

    return (
        <AppContext.Provider value={{ numberType, fieldBase }}>
            <div id="particles-js" style={{
                position: 'fixed',
                top: 0,
                left: 0,
                width: '100%',
                height: '100%',
                zIndex: 0,
                background: 'linear-gradient(to bottom right, #f8f9fa, #e9ecef)'
            }}></div>
            <div className="container mt-5 mb-5" style={{ position: 'relative', zIndex: 1 }}>
                <h1 className="page-title">
                    <i className="bi bi-matrix-fill me-2"></i>
                    MatrixGo Web
                </h1>
                <h2 className="text-center mb-4">
                    <i className={`bi ${getOperationIcon()} me-2`}></i>
                    {getOperationTitle()}
                </h2>
                
                <div className="row mb-4">
                    <div className="col-md-6 offset-md-3">
                        <div className="number-type-selector">
                            <label className="form-label">Тип чисел:</label>
                            <div className="btn-group w-100">
                                {numberTypes.map(type => (
                                    <button
                                        key={type.value}
                                        className={`btn ${numberType === type.value ? 'btn-primary' : 'btn-outline-primary'}`}
                                        onClick={() => handleNumberTypeChange(type.value)}
                                        style={{ flex: '1 1 0' }}
                                    >
                                        {type.icon} {type.label}
                                    </button>
                                ))}
                            </div>
                            {numberType === 'gf' && (
                                <div className="mt-3">
                                    <label className="form-label">Основание поля (простое число):</label>
                                    <input
                                        type="number"
                                        className={`form-control ${fieldBaseError ? 'is-invalid' : ''}`}
                                        value={fieldBaseInput}
                                        onChange={handleFieldBaseChange}
                                        min="2"
                                        placeholder="Введите простое число"
                                        style={numberInputStyle}
                                    />
                                    {fieldBaseError && (
                                        <div className="invalid-feedback">
                                            {fieldBaseError}
                                        </div>
                                    )}
                                </div>
                            )}
                        </div>
                    </div>
                </div>

                <div className="row mb-4">
                    <div className="col-md-6 offset-md-3">
                        <OperationSelect 
                            value={operation}
                            onChange={(value) => {
                                setOperation(value);
                                if (['solve', 'determinant', 'inverse'].includes(value)) {
                                    handleSizeChange(rows);
                                }
                                setResult(null);
                                setError(null);
                            }}
                        />
                    </div>
                </div>

                <div className="row mb-4">
                    <div className="col-md-6 offset-md-3">
                        <div className="input-group">
                            <span className="input-group-text">
                                <i className="bi bi-grid-3x3"></i>
                                {needsSquareMatrix ? "Размер" : "Строки"}
                            </span>
                            <SizeInput
                                value={rows}
                                onChange={(value) => handleSizeChange(value)}
                                placeholder="Размер матрицы"
                            />
                            {!needsSquareMatrix && (
                                <React.Fragment>
                                    <span className="input-group-text">×</span>
                                    <SizeInput
                                        value={cols}
                                        onChange={(value) => handleColsChange(value)}
                                        placeholder="Столбцы"
                                    />
                                </React.Fragment>
                            )}
                        </div>
                    </div>
                </div>

                <div className="row">
                    <div className="col-12">
                        <div className="matrix-equations">
                            <div className="matrix-column">
                                <h3 className="matrix-label">
                                    <i className="bi bi-grid-3x3-gap-fill me-2"></i>
                                    Матрица A
                                </h3>
                                <div className="matrix-wrapper">
                                    <Matrix 
                                        rows={rows}
                                        cols={cols}
                                        values={matrixA}
                                        onChange={handleMatrixChange}
                                    />
                                </div>
                            </div>
                            {operation === 'solve' && (
                                <div className="matrix-column">
                                    <h3 className="matrix-label">
                                        <i className="bi bi-list-ol me-2"></i>
                                        Вектор b
                                    </h3>
                                    <div className="matrix-wrapper">
                                        <Vector
                                            size={rows}
                                            values={vectorB}
                                            onChange={handleVectorChange}
                                        />
                                    </div>
                                </div>
                            )}
                        </div>
                    </div>
                </div>

                <div className="row mt-4">
                    <div className="col text-center">
                        <button 
                            className="btn btn-primary btn-lg"
                            onClick={handleSubmit}
                            disabled={loading}
                        >
                            {loading ? (
                                <span className="d-flex align-items-center justify-content-center">
                                    <span className="spinner-border spinner-border-sm" />
                                    <span className="ms-2">Вычисление...</span>
                                </span>
                            ) : (
                                <span className="d-flex align-items-center justify-content-center">
                                    <i className={`bi ${getOperationIcon()} me-2`}></i>
                                    <span>Вычислить</span>
                                </span>
                            )}
                        </button>
                    </div>
                </div>

                {error && (
                    <div className="row mt-3">
                        <div className="col">
                            <div className="error-message">
                                <i className="bi bi-exclamation-triangle-fill me-2"></i>
                                {error}
                            </div>
                        </div>
                    </div>
                )}

                {result && (
                    <div className="row mt-3">
                        <div className="col">
                            <div className="result-container">
                                <h4 className="text-center mb-4">
                                    <i className="bi bi-check-circle-fill me-2"></i>
                                    Результат:
                                </h4>
                                {typeof result === 'object' && !Array.isArray(result) ? (
                                    <p className="text-center fs-4">{formatNumber(result)}</p>
                                ) : Array.isArray(result) ? (
                                    Array.isArray(result[0]) ? (
                                        <Matrix
                                            rows={result.length}
                                            cols={result[0].length}
                                            values={result}
                                            onChange={() => {}}
                                            readOnly={true}
                                        />
                                    ) : (
                                        <Vector
                                            size={result.length}
                                            values={result}
                                            onChange={() => {}}
                                            readOnly={true}
                                        />
                                    )
                                ) : (
                                    <p className="text-center fs-4">{formatNumber({ type: 'float', value: result })}</p>
                                )}
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </AppContext.Provider>
    );
}

try {
    console.log('Starting to render App...');
    ReactDOM.render(
        <React.StrictMode>
            <App />
        </React.StrictMode>,
        document.getElementById('root'),
        () => {
            console.log('App rendered successfully');
        }
    );
} catch (error) {
    console.error('Error rendering App:', error);
    document.getElementById('root').innerHTML = `
        <div class="alert alert-danger m-5">
            Error rendering application: ${error.message}
        </div>
    `;
}

