import os

logPath = '../log/'

def read_file(filename):
    data = []
    with open(filename, 'r') as f:
        data = f.read().splitlines()
    
    return data

def write_results(results, filename = 'results.txt'):
    data = []
    haveFile = filename.split('/')[-1] in os.listdir(logPath)
    if haveFile:
        data = read_file(filename)
    data += results
    with open(filename, 'w') as f:
        for line in data:
            f.write(line + '\n')

if __name__ == '__main__':
    files = os.listdir(logPath)
    results = []
    for x in files:
        if '.log' in x:
            data = read_file(logPath + x)
            for line in data:
                if 'Max profit' in line:
                    print(line)
                    results.append(line)
    
    write_results(results, logPath + 'results.txt')