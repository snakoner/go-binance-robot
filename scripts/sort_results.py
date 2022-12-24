if __name__ == '__main__':
    data = []
    with open('../log/results.txt', 'r') as f:
        data = f.read().splitlines()
    
    data = [float(x.split(' ')[-1]) for x in data]
    data = sorted(data)
    for x in data:
        print(x)